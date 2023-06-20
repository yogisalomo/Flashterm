# Flashterm

Flash Card for your Terminal.

https://github.com/yogisalomo/flashterm/assets/4602804/8ad6baaf-eb5a-4a72-adcf-bd69ae7b9921

### Dev Environment

- Go v1.18

### Dependencies

- [PTerm](https://docs.pterm.sh/) - TUI Framework

### How to run

```shell
go run main.go
```

### TODO

- [x] Add weighted random to the next question selection.
- [ ] Generate executable so that non-programmer can use
- [ ] Create a better persistent data storage (redis, sqlite, PostgreSQL) and make it configurable
- [ ] Edit & View Vocabulary
- [ ] Fine tune the weight mechanism to make it more useful
- [ ] Code refactoring

### References

- Liu, Z. (2019, February 26). Weighted Random Sampling. Retrieved from https://zliu.org/post/weighted-random/
