package main

import (
	"errors"
	"sync"
)

type ServidorJogo struct {
	Jogadores map[string]*Jogador
	Mutex sync.Mutex
	Historico map[string]int
}

func (s *ServidorJogo) RegistrarJogador(args *Jogador, reply *bool) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.Jogadores[args.Nome] = args
	return nil
}

func (s *ServidorJogo) EnviarMensagem(args *Mensagem, reply *bool) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	for _, jogador := range s.Jogadores {
		if jogador.Nome != args.Remetente {
			// TODO: enviar a mensagem ao jogador
		}
	}

	if reply != nil {
		*reply = true
	}
	return nil
}

func (s *ServidorJogo) AtualizarPosicao(args *Movimento, reply *bool) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	jogador, ok := s.Jogadores[args.Nome]
	if !ok {
		return errors.New("Jogador não encontrado")
	}
	jogador.X = args.X
	jogador.Y = args.Y
	return nil
}

func (s *ServidorJogo) ObterEstado(args *string, reply *EstadoJogo) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	for _, jogador := range s.Jogadores {
		estado := EstadoJogo{
			Nome: jogador.Nome,
			X:    jogador.X,
			Y:    jogador.Y,
		}
		reply.Estados = append(reply.Estados, estado)
	}
	return nil
}

