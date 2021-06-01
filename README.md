# Concurrent Downloading

This project intends to reap the full benefit of goroutines, threads managed by Go runtime.

Goroutine allows the downloader to download different parts of the file, while being synced by the WaitGroup.
Once all goroutines finished its execution, different parts of the file will be merged into a single file—the original file.

_The project is still unfinished, the concurrent download functionality has yet to be implemented._
