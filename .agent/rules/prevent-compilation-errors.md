---
trigger: always_on
---

- always run 'make all' after changes to make sure that it's not broken.
- Ensure decimal values in Go test struct literals use `decimal.NewFromFloat()` instead of raw float constants to avoid type mismatch errors.
