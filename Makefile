# Makefile para o jogo multiplayer

# Compilador Go
GO = go

# Nomes dos executáveis
SERVIDOR = servidor
CLIENTE = cliente
JOGO_ORIGINAL = jogo_original

# Flags de compilação
LDFLAGS = -ldflags "-s -w"

# Regra padrão
all: servidor cliente

# Compila o servidor
servidor:
	$(GO) build $(LDFLAGS) -o $(SERVIDOR) servidor.go tipos.go

# Compila o cliente
cliente:
	$(GO) build $(LDFLAGS) -o $(CLIENTE) cliente.go jogo.go personagem.go interface.go tipos.go

# Compila o jogo original (single player)
jogo_original:
	$(GO) build $(LDFLAGS) -o $(JOGO_ORIGINAL) main.go jogo.go personagem.go interface.go tipos.go

# Executa o servidor
run_servidor: servidor
	./$(SERVIDOR)

# Executa o cliente (requer nome do jogador)
run_cliente: cliente
	@if [ -z "$(NOME)" ]; then \
		echo "Uso: make run_cliente NOME=nome_do_jogador"; \
		exit 1; \
	fi
	./$(CLIENTE) $(NOME)

# Executa o jogo original
run_original: jogo_original
	./$(JOGO_ORIGINAL)

# Limpa arquivos compilados
clean:
	rm -f $(SERVIDOR) $(CLIENTE) $(JOGO_ORIGINAL)

# Instala dependências
deps:
	$(GO) mod init fppd-jogo-t2
	$(GO) get github.com/nsf/termbox-go

# Ajuda
help:
	@echo "Comandos disponíveis:"
	@echo "  make all          - Compila servidor e cliente"
	@echo "  make servidor     - Compila apenas o servidor"
	@echo "  make cliente      - Compila apenas o cliente"
	@echo "  make run_servidor - Executa o servidor"
	@echo "  make run_cliente NOME=nome - Executa cliente com nome específico"
	@echo "  make run_original - Executa jogo original (single player)"
	@echo "  make clean        - Remove arquivos compilados"
	@echo "  make deps         - Instala dependências"
	@echo ""
	@echo "Exemplo de uso:"
	@echo "  1. make run_servidor (em um terminal)"
	@echo "  2. make run_cliente NOME=Joao (em outro terminal)"
	@echo "  3. make run_cliente NOME=Maria (em outro terminal)"

.PHONY: all servidor cliente jogo_original run_servidor run_cliente run_original clean deps help