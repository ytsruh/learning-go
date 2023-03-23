# Hacker News

The original client was already provided and the goal was to improve it as much as possible

## Order

1. Original
2. First Attempt = Added a single blocking go routine that made no improvement to load time
3. Concurrency = Added multiple go routines to speed up the load times
4. 30 Stories = Increased go routine loop to ensure that 30 stories are returned
5. Cache Mutex = Added simple cache to increase speed. Mutex also added to avoid race condition
6. Better Cache = Added a background go routine to update the cache before it expires
