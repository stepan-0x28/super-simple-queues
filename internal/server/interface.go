package server

type Server interface {
	Run(port int) error
}

func RunGo(server Server, port int, errChan chan error) {
	go func() {
		if err := server.Run(port); err != nil {
			errChan <- err
		}
	}()
}
