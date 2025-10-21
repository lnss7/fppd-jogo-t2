// main.go - Loop principal do jogo
package main

import (
	"fmt"
	"net/rpc"
	"os"
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

	cliente, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		panic ("Erro ao conectar no servidor: " + err.Error())
	}

	var nome string 
	fmt.Print("Digite seu nome:")
	fmt.Scanln(&nome)

	jogador := Jogador{Nome: nome, X: jogo.PosX, Y: jogo.PosY}
	var ok bool 
	cliente.Call("Servidor.Jogo.RegistrarJogador", jogador, &ok)

	// Desenha o estado inicial do jogo
	interfaceDesenharJogo(&jogo)

	// Loop principal de entrada
	for {
		evento := interfaceLerEventoTeclado()
		if continuar := personagemExecutarAcao(evento, &jogo); !continuar {
			break
		}

		var estado EstadoJogo
		cliente.Call("ServidorJogo.ObterEstado", &nome, &estado)

		for _, j := range estado.Estados {
			if j.Nome != nome {
				jogo.Mapa[j.Y][j.X] = Personagem
			}
		}

		fmt.Println("todos os players:")
		for _, u := range estado.Estados {
        fmt.Printf("Nome: %d, X: %s, Y: %s\n", u.Nome, u.X, u.Y)
    }
		interfaceDesenharJogo(&jogo)
	}
}