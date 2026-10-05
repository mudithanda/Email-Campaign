# Concurrent Email Campaign Sender

A tool written in Go that sends personalized emails to a list of recipients using a **concurrent worker pool**. It reads recipients from a CSV file, renders a personalized message for each one from an email template, and delivers them over SMTP. For local development it works with [Mailpit](https://github.com/axllent/mailpit), a fake SMTP server with a web inbox, so no real emails are sent while testing.

## What this project is about

Sending thousands of emails one at a time is slow. This project shows how to speed that up with Go's concurrency primitives while keeping the code simple:

- One **producer goroutine** reads the CSV and pushes each recipient onto a channel.
- Several **worker goroutines** (5 by default) pull recipients from that channel, render the template, and send the email in parallel.
- A sync.WaitGroup keeps main alive until every worker has finished.

It is a small, practical example of the producer/consumer pattern, channels, goroutines, and graceful shutdown by closing a channel.

## How it works

```
emails.csv ──► loadRecipient (producer) ──► unbuffered channel ──┬─► worker 1 ─┐
                                                                 ├─► worker 2 ─┤
                                                                 ├─► worker 3 ─┼─► SMTP server
                                                                 ├─► worker 4 ─┤   (Mailpit)
                                                                 └─► worker 5 ─┘
```

1. `main.go` creates an unbuffered `chan Recipient`, starts the producer goroutine, and launches the worker pool.
2. `loadRecipient` (in the CSV loader file) opens the CSV, skips the header row, trims whitespace, skips malformed rows, and sends each `Recipient{Name, Email}` into the channel. When it finishes, it closes the channel.
3. Each `emailWorker` loops over the channel (`for recipient := range ch`). For every recipient it:
   - renders `email.tmpl` with the recipient's data (`executeTemplate`),
   - sends the result with `net/smtp`,
   - waits 50 ms to avoid overwhelming the SMTP server,
   - logs success or failure.
4. When the channel is closed and drained, the workers exit, `wg.Wait()` returns, and the program ends.

## Project structure

| File | Purpose |
|------|---------|
| `main.go` | Entry point: creates the channel, starts the producer and workers, defines `Recipient` and `executeTemplate` |
| `producer.go` (CSV loader) | `loadRecipient`: reads and validates the CSV, then feeds recipients into the channel |
| `consumer.go` (email worker) | `emailWorker`: renders templates and sends emails over SMTP |
| `email.tmpl` | Email template (Go template syntax) |
| `emails.csv` | Recipient list (`name,email`) |
| `info.md` | Docker command to run the Mailpit SMTP server |

Adjust the file names above if yours differ.

## Prerequisites

- [Go](https://go.dev/dl/) 1.18 or later
- [Mailpit](https://github.com/axllent/mailpit) (or any SMTP server listening on `localhost:1025`)
- Docker, if you want to run Mailpit as a container

## Setup

### 1. Start Mailpit

Run the command from `info.md`. A typical Docker command is:

```bash
docker run -d --name mailpit -p 8025:8025 -p 1025:1025 axllent/mailpit
```

- SMTP server: `localhost:1025`
- Web inbox: http://localhost:8025

### 2. Prepare the recipient list

`emails.csv` must have a header row followed by one recipient per line, in the order `name,email`:

```csv
name,email
User1,user1@gmail.com
User2,user2@gmail.com
```

Rows with fewer than two columns are skipped with a log message.

### 3. Edit the email template

`email.tmpl` is a standard Go template. The recipient's fields are available as `{{.Name}}` and `{{.Email}}`. Because the whole template output is passed to `smtp.SendMail` as the raw message, include the headers at the top:

```
To: {{.Email}}
Subject: Hello {{.Name}}

Hi {{.Name}}

Thanks,
The MH Campaign Team.
```

Keep a blank line between the headers and the body.

### 4. Run

```bash
go run .
```

Then open http://localhost:8025 to see the delivered emails in Mailpit.

Example console output:

```
Worker 3: Sending Email to "user1@gmail.com"
Worker 1: Sending Email to "user2@gmail.com"
Worker 3: Sent Email to user1@gmail.com
Worker 1: Sent Email to user2@gmail.com
```

## Concepts demonstrated

- Goroutines and the producer/consumer pattern
- Unbuffered channels for hand-off between goroutines
- Worker pool with `sync.WaitGroup`
- Closing a channel to signal completion, with `range` over a channel
- CSV parsing with `encoding/csv`
- Template rendering with the standard library
- Sending mail with `net/smtp`

