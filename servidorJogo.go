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

	// adicionado: comandos já processados por jogador (seq -> true)
    Processed map[string]map[int]bool
}

func main() {
    servidor := &ServidorJogo{
        Jogadores: make(map[string]*Jogador),
        Historico: make(map[string]int),
		Processed: make(map[string]map[int]bool),
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

func (s *ServidorJogo) RegistrarJogador(args *CmdJogador, reply *bool) error {
    s.Mutex.Lock()

    nome := args.Jogador.Nome
    if s.Processed[nome] == nil {
        s.Processed[nome] = make(map[int]bool)
    }
    if s.Processed[nome][args.Seq] {
        if reply != nil {
            *reply = true
        }
        // já processado -> idempotente
		s.Mutex.Unlock()
        return nil
    }

    // processa comando
    s.Jogadores[nome] = &args.Jogador
    s.Historico[nome] = int(time.Now().Unix())
    s.Processed[nome][args.Seq] = true

    if reply != nil {
        *reply = true
    }
    fmt.Printf("Jogador registrado: %s (%d,%d)\n", args.Jogador.Nome, args.Jogador.X, args.Jogador.Y)
    // unlock antes de PrintEstado
    s.Mutex.Unlock()
    s.PrintEstado()
    return nil
}

func (s *ServidorJogo) EnviarMensagem(args *Mensagem, reply *bool) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	for _, jogador := range s.Jogadores {
		if jogador.Nome != args.Remetente {
			// TODO: enviar a mnsagem ao jogador
		}
	}

	if reply != nil {
		*reply = true
	}
	return nil
}

func (s *ServidorJogo) AtualizarPosicao(args *CmdMovimento, reply *bool) error {
    s.Mutex.Lock()

    nome := args.Movimento.Nome
    if s.Processed[nome] == nil {
        s.Processed[nome] = make(map[int]bool)
    }
    if s.Processed[nome][args.Seq] {
        s.Mutex.Unlock()
        if reply != nil {
            *reply = true
        }
        return nil
    }

    jogador, ok := s.Jogadores[nome]
    if !ok {
        s.Mutex.Unlock()
        return errors.New("Jogador não encontrado")
    }
    jogador.X = args.Movimento.X
    jogador.Y = args.Movimento.Y
    s.Historico[nome] = int(time.Now().Unix())
    s.Processed[nome][args.Seq] = true

    if reply != nil {
        *reply = true
    }
    fmt.Printf("Posição atualizada: %s -> (%d,%d)\n", nome, args.Movimento.X, args.Movimento.Y)
    s.Mutex.Unlock()
    s.PrintEstado()
    return nil
}
func (s *ServidorJogo) RemoverJogador(args *CmdRemover, reply *bool) error {
    s.Mutex.Lock()

    nome := args.Nome
    if s.Processed[nome] == nil {
        s.Processed[nome] = make(map[int]bool)
    }
    if s.Processed[nome][args.Seq] {
        s.Mutex.Unlock()
        if reply != nil {
            *reply = true
        }
        return nil
    }

    if _, ok := s.Jogadores[nome]; ok {
        delete(s.Jogadores, nome)
        delete(s.Historico, nome)
        delete(s.Processed, nome)
        s.Processed[nome][args.Seq] = true
        if reply != nil {
            *reply = true
        }
        fmt.Printf("Jogador removido: %s\n", nome)
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

