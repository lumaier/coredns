package plugin

func PrintMemUsage() {
	// // Force GC to clear up, should see a memory drop
	// runtime.GC()

	// var m runtime.MemStats
	// runtime.ReadMemStats(&m)
	// // For info on each, see: https://golang.org/pkg/runtime/#MemStats
	// fmt.Printf("Alloc = %v MB", bToMb(m.Alloc))
	// fmt.Printf("\tTotalAlloc = %v MB", bToMb(m.TotalAlloc))
	// fmt.Printf("\tStackSys = %v MB", bToMb(m.StackSys))
	// fmt.Printf("\tSys = %v MB", bToMb(m.Sys))
	// fmt.Printf("\tNumGC = %v\n", m.NumGC)
	// os.Exit(3)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
