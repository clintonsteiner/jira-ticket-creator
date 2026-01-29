go run main.go \
  -url https://your-domain.atlassian.net \
  -email you@example.com \
  -token your-api-token \
  -project PROJECT1 \
  "Implement login feature" \
  "Login feature blocked until core setup is done" \
  "PROJECT1-12,PROJECT1-15"

go run main.go report

