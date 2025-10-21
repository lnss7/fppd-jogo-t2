package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"
)

type ServidorJogo struct {
	Jogadores        map[string]*Jogador
	Mutex            sync.RWMutex
	SequenceNumbers  map[string]int // Controle de sequence numbers por jogador
	ComandosProcessados map[string]map[int]bool // Comandos já processados por jogador
}

func main() {
	mainServidor()
}

func mainServidor() {
	servidor := &ServidorJogo{
		Jogadores:        make(map[string]*Jogador),
		SequenceNumbers:  make(map[string]int),
		ComandosProcessados: make(map[string]map[int]bool),
	}

	// Registra os métodos RPC
	rpc.Register(servidor)

	// Inicia o servidor na porta 1234
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("Erro ao iniciar servidor:", err)
	}

	fmt.Println("Servidor RPC rodando na porta 1234...")
	fmt.Println("Aguardando conexões de clientes...")
	
	// Aceita conexões indefinidamente
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Erro ao aceitar conexão: %v", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}

// RegistrarJogador registra um novo jogador no servidor
func (s *ServidorJogo) RegistrarJogador(args *RegistroJogador, reply *bool) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	fmt.Printf("Cliente solicitando registro: %s\n", args.Nome)

	// Verifica se o jogador já existe
	if _, existe := s.Jogadores[args.Nome]; existe {
		fmt.Printf("Jogador %s já está registrado\n", args.Nome)
		*reply = false
		return nil
	}

	// Cria novo jogador
	s.Jogadores[args.Nome] = &Jogador{
		Nome:              args.Nome,
		X:                 args.X,
		Y:                 args.Y,
		UltimaAtualizacao: time.Now(),
	}

	// Inicializa controle de sequence numbers
	s.SequenceNumbers[args.Nome] = 0
	s.ComandosProcessados[args.Nome] = make(map[int]bool)

	fmt.Printf("Jogador %s registrado com sucesso na posição (%d, %d)\n", args.Nome, args.X, args.Y)
	*reply = true
	return nil
}

// ProcessarComando processa comandos dos jogadores com garantia de execução única
func (s *ServidorJogo) ProcessarComando(args *Comando, reply *RespostaComando) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	fmt.Printf("Processando comando de %s: %s (seq: %d)\n", args.Nome, args.Tipo, args.SequenceNumber)

	// Verifica se o jogador existe
	jogador, existe := s.Jogadores[args.Nome]
	if !existe {
		reply.Sucesso = false
		reply.Mensagem = "Jogador não encontrado"
		return nil
	}

	// Verifica se o comando já foi processado (garantia de execução única)
	if s.ComandosProcessados[args.Nome][args.SequenceNumber] {
		fmt.Printf("Comando %d de %s já foi processado, ignorando\n", args.SequenceNumber, args.Nome)
		reply.Sucesso = true
		reply.Mensagem = "Comando já processado"
		return nil
	}

	// Verifica se o sequence number é válido (deve ser sequencial)
	expectedSeq := s.SequenceNumbers[args.Nome] + 1
	if args.SequenceNumber != expectedSeq {
		fmt.Printf("Sequence number inválido para %s: esperado %d, recebido %d\n", 
			args.Nome, expectedSeq, args.SequenceNumber)
		reply.Sucesso = false
		reply.Mensagem = "Sequence number inválido"
		return nil
	}

	// Processa o comando baseado no tipo
	switch args.Tipo {
	case "mover":
		jogador.X = args.X
		jogador.Y = args.Y
		jogador.UltimaAtualizacao = time.Now()
		fmt.Printf("Jogador %s moveu para (%d, %d)\n", args.Nome, args.X, args.Y)
	case "interagir":
		fmt.Printf("Jogador %s interagiu na posição (%d, %d)\n", args.Nome, args.X, args.Y)
		// Aqui você pode adicionar lógica de interação específica
	default:
		reply.Sucesso = false
		reply.Mensagem = "Tipo de comando inválido"
		return nil
	}

	// Marca o comando como processado
	s.ComandosProcessados[args.Nome][args.SequenceNumber] = true
	s.SequenceNumbers[args.Nome] = args.SequenceNumber

	reply.Sucesso = true
	reply.Mensagem = "Comando processado com sucesso"
	return nil
}

// ObterEstado retorna o estado atual de todos os jogadores
func (s *ServidorJogo) ObterEstado(args *string, reply *EstadoJogo) error {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	fmt.Printf("Cliente %s solicitando estado do jogo\n", *args)

	// Cria lista de estados dos jogadores
	reply.Estados = make([]Jogador, 0, len(s.Jogadores))
	for _, jogador := range s.Jogadores {
		reply.Estados = append(reply.Estados, *jogador)
	}

	fmt.Printf("Retornando estado de %d jogadores\n", len(reply.Estados))
	return nil
}

// DesregistrarJogador remove um jogador do servidor
func (s *ServidorJogo) DesregistrarJogador(args *string, reply *bool) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	fmt.Printf("Desregistrando jogador: %s\n", *args)

	if _, existe := s.Jogadores[*args]; !existe {
		*reply = false
		return nil
	}

	delete(s.Jogadores, *args)
	delete(s.SequenceNumbers, *args)
	delete(s.ComandosProcessados, *args)

	fmt.Printf("Jogador %s desregistrado com sucesso\n", *args)
	*reply = true
	return nil
}

// Ping verifica se o servidor está ativo
func (s *ServidorJogo) Ping(args *string, reply *bool) error {
	*reply = true
	return nil
}
