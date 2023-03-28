package server

func (s *Server) InitRoutes() {
	s.router.POST("/register", s.register.Handle())
	s.router.POST("/finish-register", s.finishRegister.Handle())
}
