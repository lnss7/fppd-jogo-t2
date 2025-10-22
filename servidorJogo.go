package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"
)

type ServidorJogo struct {
	Jogadores map[string]*Jogador
	Mutex sync.Mutex
	Historico map[string]int
}

func main() {
    servidor := &ServidorJogo{
        Jogadores: make(map[string]*Jogador),
        Historico: make(map[string]int),
    }

    rpc.Register(servidor)

	// goroutine que imprime periodicamente o estado do servidor
	go func() {
		for {
			servidor.PrintEstado()
			time.Sleep(5 * time.Second)
		}
	}()

    listener, err := net.Listen("tcp", ":1234")
    if err != nil {
        log.Fatal("Erro ao iniciar servidor:", err)
    }

    fmt.Println("Servidor RPC rodando na porta 1234...")
    rpc.Accept(listener)
}

func (s *ServidorJogo) RegistrarJogador(args *Jogador, reply *bool) error {
	s.Mutex.Lock()
	s.Jogadores[args.Nome] = args

	// registra timestamp de último contato
    s.Historico[args.Nome] = int(time.Now().Unix())

	if reply != nil {
		*reply = true
	}
	fmt.Printf("Jogador registrado: %s (%d,%d)\n", args.Nome, args.X, args.Y)
	s.Mutex.Unlock()
	s.PrintEstado()
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

	jogador, ok := s.Jogadores[args.Nome]
	if !ok {
		s.Mutex.Unlock() // desbloquia antes pra nao dar deadlock
		return errors.New("Jogador não encontrado")
	}
	jogador.X = args.X
	jogador.Y = args.Y

    // atualiza timestamp de último contato
    s.Historico[args.Nome] = int(time.Now().Unix())

	if reply != nil {
		*reply = true
	}
	fmt.Printf("Posição atualizada: %s -> (%d,%d)\n", args.Nome, args.X, args.Y)
	s.Mutex.Unlock()
	s.PrintEstado()
	return nil
}

func (s *ServidorJogo) RemoverJogador(args *string, reply *bool) error {
    s.Mutex.Lock() 

    if _, ok := s.Jogadores[*args]; ok {
        delete(s.Jogadores, *args)
        delete(s.Historico, *args)
        if reply != nil {
            *reply = true
        }
        fmt.Printf("Jogador removido: %s\n", *args)
		s.Mutex.Unlock()
        s.PrintEstado()
        return nil
    }

	s.Mutex.Unlock()
    if reply != nil {
        *reply = false
    }
    return errors.New("Jogador não encontrado")
}

// PrintEstado escreve o estado atual dos jogadores registrados.
func (s *ServidorJogo) PrintEstado() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	fmt.Println("Estado do Servidor: jogadores registrados")
	if len(s.Jogadores) == 0 {
		fmt.Println("nenhum jogador conectado")
		fmt.Println("--------------------------------------------------")
		return
	}
	for _, j := range s.Jogadores {
		fmt.Printf("Nome: %s, X: %d, Y: %d\n", j.Nome, j.X, j.Y)
	}
	fmt.Println("--------------------------------------------------")
}

func (s *ServidorJogo) ObterEstado(args *string, reply *EstadoJogo) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	if reply == nil {
		return errors.New("reply is nil")
	}

	if reply.Estados == nil {
		reply.Estados = make([]Jogador, 0, len(s.Jogadores))
	}

	for _, jogador := range s.Jogadores {
		reply.Estados = append(reply.Estados, *jogador)
	}

	return nil
}

