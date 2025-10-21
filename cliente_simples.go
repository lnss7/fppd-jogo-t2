package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"sync"
	"time"
)

type ClienteJogo struct {
	Nome           string
	PosX, PosY     int
	SequenceNumber int
	ClienteRPC     *rpc.Client
	Jogo           *Jogo
	OutrosJogadores map[string]Jogador
	Mutex          sync.RWMutex
}

var cliente *ClienteJogo

func main() {
	// Verifica argumentos
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run cliente_simples.go <nome_do_jogador> [endereco_servidor]")
		fmt.Println("Exemplo: go run cliente_simples.go Joao localhost:1234")
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

	// Inicializa a interface
	interfaceIniciar()
	defer interfaceFinalizar()

	// Carrega o mapa
	mapaFile := "mapa.txt"
	if len(os.Args) > 3 {
		mapaFile = os.Args[3]
	}

	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	// Cria cliente
	cliente = &ClienteJogo{
		Nome:            nomeJogador,
		PosX:            jogo.PosX,
		PosY:            jogo.PosY,
		SequenceNumber:  0,
		ClienteRPC:      clienteRPC,
		Jogo:            &jogo,
		OutrosJogadores: make(map[string]Jogador),
	}

	// Registra no servidor
	if err := cliente.registrarNoServidor(); err != nil {
		log.Fatal("Erro ao registrar no servidor:", err)
	}

	// Inicia thread para buscar atualizações
	go cliente.buscarAtualizacoes()

	// Desenha estado inicial
	cliente.desenharJogo()

	// Loop principal de entrada
	for {
		evento := interfaceLerEventoTeclado()
		if continuar := cliente.processarEvento(evento); !continuar {
			break
		}
		cliente.desenharJogo()
	}

	// Desregistra do servidor
	cliente.desregistrarDoServidor()
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

	fmt.Printf("Jogador %s registrado com sucesso!\n", c.Nome)
	return nil
}

func (c *ClienteJogo) desregistrarDoServidor() {
	var sucesso bool
	c.ClienteRPC.Call("ServidorJogo.DesregistrarJogador", &c.Nome, &sucesso)
	if sucesso {
		fmt.Printf("Jogador %s desregistrado\n", c.Nome)
	}
}

func (c *ClienteJogo) buscarAtualizacoes() {
	for {
		time.Sleep(500 * time.Millisecond) // Busca atualizações a cada 500ms

		var estado EstadoJogo
		err := c.ClienteRPC.Call("ServidorJogo.ObterEstado", &c.Nome, &estado)
		if err != nil {
			continue // Ignora erros de rede temporários
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

func (c *ClienteJogo) processarEvento(ev EventoTeclado) bool {
	switch ev.Tipo {
	case "sair":
		return false
	case "interagir":
		c.enviarComando("interagir", c.PosX, c.PosY)
	case "mover":
		dx, dy := 0, 0
		switch ev.Tecla {
		case 'w': dy = -1
		case 'a': dx = -1
		case 's': dy = 1
		case 'd': dx = 1
		}

		nx, ny := c.PosX+dx, c.PosY+dy
		if c.jogoPodeMoverPara(nx, ny) {
			c.enviarComando("mover", nx, ny)
			c.PosX, c.PosY = nx, ny
			c.Jogo.PosX, c.Jogo.PosY = nx, ny
		}
	}
	return true
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

func (c *ClienteJogo) jogoPodeMoverPara(x, y int) bool {
	return jogoPodeMoverPara(c.Jogo, x, y)
}

func (c *ClienteJogo) desenharJogo() {
	interfaceLimparTela()

	// Desenha o mapa
	for y, linha := range c.Jogo.Mapa {
		for x, elem := range linha {
			interfaceDesenharElemento(x, y, elem)
		}
	}

	// Desenha outros jogadores
	c.Mutex.RLock()
	for _, jogador := range c.OutrosJogadores {
		interfaceDesenharElemento(jogador.X, jogador.Y, Inimigo)
	}
	c.Mutex.RUnlock()

	// Desenha o jogador atual
	interfaceDesenharElemento(c.PosX, c.PosY, Personagem)

	// Desenha barra de status
	c.desenharBarraDeStatus()

	interfaceAtualizarTela()
}

func (c *ClienteJogo) desenharBarraDeStatus() {
	// Status do jogador atual
	status := fmt.Sprintf("Jogador: %s | Pos: (%d,%d) | Seq: %d", 
		c.Nome, c.PosX, c.PosY, c.SequenceNumber)
	for i, ch := range status {
		interfaceDesenharElemento(i, len(c.Jogo.Mapa)+1, Elemento{ch, CorTexto, CorPadrao, false})
	}

	// Lista de outros jogadores
	c.Mutex.RLock()
	linha := 2
	for nome, jogador := range c.OutrosJogadores {
		info := fmt.Sprintf("%s: (%d,%d)", nome, jogador.X, jogador.Y)
		for i, ch := range info {
			interfaceDesenharElemento(i, len(c.Jogo.Mapa)+linha, Elemento{ch, CorTexto, CorPadrao, false})
		}
		linha++
	}
	c.Mutex.RUnlock()

	// Instruções
	instrucoes := "WASD=mover E=interagir ESC=sair"
	for i, ch := range instrucoes {
		interfaceDesenharElemento(i, len(c.Jogo.Mapa)+linha+1, Elemento{ch, CorTexto, CorPadrao, false})
	}
}
