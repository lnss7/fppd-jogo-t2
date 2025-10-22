// main.go - Loop principal do jogo
package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"sync"
	"time"
)

func main() {
	// Inicializa a interface (termbox)
	interfaceIniciar()
	defer interfaceFinalizar()

	// Usa "mapa.txt" como arquivo padrão ou lê o primeiro argumento
	mapaFile := "mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}

	// Inicializa o jogo
	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	clienteAddr := "localhost:1234"
    // se passado como 2º argumento (ex: mapa.txt 192.168.1.10:1234)
    if len(os.Args) > 2 {
        clienteAddr = os.Args[2]
    }

    // loga os argumentos e o endereço a ser usado (ajuda debugging)
    log.Printf("os.Args = %v", os.Args)
    log.Printf("Endereço do servidor: %s", clienteAddr)

    // tenta conectar com retry antes de falhar
    var cliente *rpc.Client
    var err error
    maxTentativas := 5
    delay := 500 * time.Millisecond
    for i := 0; i < maxTentativas; i++ {
        cliente, err = rpc.Dial("tcp", clienteAddr)
        if err == nil {
            break
        }
        log.Printf("Falha ao conectar (%d/%d) em %s: %v", i+1, maxTentativas, clienteAddr, err)
        time.Sleep(delay)
        delay *= 2
    }
    if err != nil {
        log.Fatalf("Não foi possível conectar no servidor %s: %v", clienteAddr, err)
    }

	interfaceFinalizar()
	var nome string
	fmt.Print("Digite seu nome: ")
	fmt.Scanln(&nome)
	interfaceIniciar()
	jogador := Jogador{Nome: nome, X: jogo.PosX, Y: jogo.PosY}
	var ok bool
	if err := cliente.Call("ServidorJogo.RegistrarJogador", &jogador, &ok); err != nil {
		log.Fatal("Erro ao registrar jogador:", err)
	}
	// Desenha o estado inicial do jogo
	interfaceDesenharJogo(&jogo)

	// mutex para proteger acesso concorrente ao jogo entre polling e loop principal
	var mu2 sync.Mutex

	// goroutine de polling: obtém estado do servidor a cada 300ms
	go func() {
		ticker := time.NewTicker(300 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			var estado EstadoJogo
			if err := cliente.Call("ServidorJogo.ObterEstado", &nome, &estado); err != nil {
				// se houver erro, apenas log e segue
				log.Println("Erro ao obter estado (poll):", err)
				continue
			}

			// atualiza mapa: remove personagens antigos (exceto o próprio) e coloca os atuais
			mu2.Lock()
			for y := range jogo.Mapa {
				for x := range jogo.Mapa[y] {
					if jogo.Mapa[y][x] == Personagem {
						if !(x == jogo.PosX && y == jogo.PosY) {
							jogo.Mapa[y][x] = Vazio
						}
					}
				}
			}
			for _, j := range estado.Estados {
				if j.Nome == nome {
					continue
				}
				if j.Y >= 0 && j.Y < len(jogo.Mapa) && j.X >= 0 && j.X < len(jogo.Mapa[j.Y]) {
					jogo.Mapa[j.Y][j.X] = Personagem
				}
			}
			interfaceDesenharJogo(&jogo)
			mu2.Unlock()
		}
	}()

	// Loop principal de entrada
	for {
		evento := interfaceLerEventoTeclado()
		if continuar := personagemExecutarAcao(evento, &jogo); !continuar {
			break
		}

		// envia a nova posição para o servidor sempre que houver movimento
		var ok bool
		mov := Movimento{Nome: nome, X: jogo.PosX, Y: jogo.PosY}
		if err := cliente.Call("ServidorJogo.AtualizarPosicao", &mov, &ok); err != nil {
			log.Println("Erro ao atualizar posicao:", err)
		}

		var estado EstadoJogo
		if err := cliente.Call("ServidorJogo.ObterEstado", &nome, &estado); err != nil {
			log.Println("Erro ao obter estado:", err)
		}

		for _, j := range estado.Estados {
			if j.Nome != nome {
				jogo.Mapa[j.Y][j.X] = Personagem
			}
		}
        // opcional: redesenhar (polling também redesenha)
        mu2.Lock()
        interfaceDesenharJogo(&jogo)
        mu2.Unlock()	
	}
	// ao sair, avisa o servidor para remover o jogador
	var removed bool
	if err := cliente.Call("ServidorJogo.RemoverJogador", &nome, &removed); err != nil {
		log.Println("Erro ao remover jogador no servidor:", err)
	}
}
