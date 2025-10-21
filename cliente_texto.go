package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strings"
	"sync"
	"time"
)

type ClienteJogo struct {
	Nome           string
	PosX, PosY     int
	SequenceNumber int
	ClienteRPC     *rpc.Client
	OutrosJogadores map[string]Jogador
	Mutex          sync.RWMutex
}

var cliente *ClienteJogo

func main() {
	// Verifica argumentos
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run cliente_texto.go <nome_do_jogador> [endereco_servidor]")
		fmt.Println("Exemplo: go run cliente_texto.go Joao localhost:1234")
		os.Exit(1)
	}

	nomeJogador := os.Args[1]
	enderecoServidor := "localhost:1234"
	if len(os.Args) > 2 {
		enderecoServidor = os.Args[2]
	}

	// Conecta ao servidor
	clienteRPC, err := rpc.Dial("tcp", enderecoServidor)
	if err != nil {
		log.Fatal("Erro ao conectar ao servidor:", err)
	}
	defer clienteRPC.Close()

	// Cria cliente
	cliente = &ClienteJogo{
		Nome:            nomeJogador,
		PosX:            5,
		PosY:            5,
		SequenceNumber:  0,
		ClienteRPC:      clienteRPC,
		OutrosJogadores: make(map[string]Jogador),
	}

	// Registra no servidor
	if err := cliente.registrarNoServidor(); err != nil {
		log.Fatal("Erro ao registrar no servidor:", err)
	}

	// Inicia thread para buscar atualizações
	go cliente.buscarAtualizacoes()

	fmt.Printf("=== JOGO MULTIPLAYER ===\n")
	fmt.Printf("Jogador: %s\n", cliente.Nome)
	fmt.Printf("Posição inicial: (%d, %d)\n", cliente.PosX, cliente.PosY)
	fmt.Printf("Comandos: w/a/s/d (mover), e (interagir), q (sair)\n\n")

	// Loop principal de entrada
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Comando: ")
		if !scanner.Scan() {
			break
		}

		comando := strings.ToLower(strings.TrimSpace(scanner.Text()))
		
		if comando == "q" {
			break
		}

		cliente.processarComando(comando)
		cliente.mostrarEstado()
	}

	// Desregistra do servidor
	cliente.desregistrarDoServidor()
	fmt.Println("Jogo encerrado!")
}

func (c *ClienteJogo) registrarNoServidor() error {
	registro := &RegistroJogador{
		Nome: c.Nome,
		X:    c.PosX,
		Y:    c.PosY,
	}

	var sucesso bool
	err := c.ClienteRPC.Call("ServidorJogo.RegistrarJogador", registro, &sucesso)
	if err != nil {
		return err
	}

	if !sucesso {
		return fmt.Errorf("falha ao registrar jogador (nome já existe?)")
	}

	return nil
}

func (c *ClienteJogo) desregistrarDoServidor() {
	var sucesso bool
	c.ClienteRPC.Call("ServidorJogo.DesregistrarJogador", &c.Nome, &sucesso)
}

func (c *ClienteJogo) buscarAtualizacoes() {
	for {
		time.Sleep(500 * time.Millisecond)

		var estado EstadoJogo
		err := c.ClienteRPC.Call("ServidorJogo.ObterEstado", &c.Nome, &estado)
		if err != nil {
			continue
		}

		// Atualiza lista de outros jogadores
		c.Mutex.Lock()
		c.OutrosJogadores = make(map[string]Jogador)
		for _, jogador := range estado.Estados {
			if jogador.Nome != c.Nome {
				c.OutrosJogadores[jogador.Nome] = jogador
			}
		}
		c.Mutex.Unlock()
	}
}

func (c *ClienteJogo) processarComando(comando string) {
	switch comando {
	case "w":
		c.mover(0, -1)
	case "a":
		c.mover(-1, 0)
	case "s":
		c.mover(0, 1)
	case "d":
		c.mover(1, 0)
	case "e":
		c.interagir()
	default:
		fmt.Println("Comando inválido! Use: w/a/s/d (mover), e (interagir), q (sair)")
	}
}

func (c *ClienteJogo) mover(dx, dy int) {
	nx, ny := c.PosX+dx, c.PosY+dy
	
	// Validação básica de limites
	if nx < 0 || ny < 0 || nx > 20 || ny > 20 {
		fmt.Println("Movimento inválido: fora dos limites")
		return
	}

	c.enviarComando("mover", nx, ny)
	c.PosX, c.PosY = nx, ny
	fmt.Printf("Movido para (%d, %d)\n", nx, ny)
}

func (c *ClienteJogo) interagir() {
	c.enviarComando("interagir", c.PosX, c.PosY)
	fmt.Printf("Interagindo na posição (%d, %d)\n", c.PosX, c.PosY)
}

func (c *ClienteJogo) enviarComando(tipo string, x, y int) {
	c.SequenceNumber++
	comando := &Comando{
		Nome:           c.Nome,
		X:              x,
		Y:              y,
		SequenceNumber: c.SequenceNumber,
		Tipo:           tipo,
	}

	var resposta RespostaComando
	err := c.ClienteRPC.Call("ServidorJogo.ProcessarComando", comando, &resposta)
	if err != nil {
		fmt.Printf("Erro ao enviar comando: %v\n", err)
		return
	}

	if !resposta.Sucesso {
		fmt.Printf("Comando rejeitado: %s\n", resposta.Mensagem)
	}
}

func (c *ClienteJogo) mostrarEstado() {
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()

	fmt.Printf("\n--- ESTADO ATUAL ---\n")
	fmt.Printf("Você: %s em (%d, %d) [seq: %d]\n", c.Nome, c.PosX, c.PosY, c.SequenceNumber)
	
	if len(c.OutrosJogadores) > 0 {
		fmt.Printf("Outros jogadores:\n")
		for nome, jogador := range c.OutrosJogadores {
			fmt.Printf("  - %s: (%d, %d)\n", nome, jogador.X, jogador.Y)
		}
	} else {
		fmt.Printf("Nenhum outro jogador conectado\n")
	}
	fmt.Printf("-------------------\n\n")
}
