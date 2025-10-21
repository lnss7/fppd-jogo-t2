#!/bin/bash

# Script para testar o jogo multiplayer
# Este script inicia o servidor e alguns clientes para demonstração

echo "=== Teste do Jogo Multiplayer ==="
echo ""

# Verifica se os executáveis existem
if [ ! -f "./servidor" ] || [ ! -f "./cliente" ]; then
    echo "Compilando executáveis..."
    make all
    echo ""
fi

echo "Iniciando servidor em background..."
./servidor &
SERVIDOR_PID=$!

# Aguarda o servidor inicializar
sleep 2

echo "Servidor iniciado (PID: $SERVIDOR_PID)"
echo ""

echo "Iniciando clientes de teste..."
echo "Pressione Ctrl+C para parar todos os processos"
echo ""

# Inicia alguns clientes de exemplo
echo "Iniciando cliente 'Jogador1'..."
gnome-terminal --title="Jogador1" -- ./cliente Jogador1 &
sleep 1

echo "Iniciando cliente 'Jogador2'..."
gnome-terminal --title="Jogador2" -- ./cliente Jogador2 &
sleep 1

echo "Iniciando cliente 'Jogador3'..."
gnome-terminal --title="Jogador3" -- ./cliente Jogador3 &

echo ""
echo "Clientes iniciados! Use WASD para mover e E para interagir."
echo "Cada jogador verá os outros como símbolos de caveira (☠)"
echo ""

# Função para limpar processos ao sair
cleanup() {
    echo ""
    echo "Parando servidor e clientes..."
    kill $SERVIDOR_PID 2>/dev/null
    pkill -f "./cliente" 2>/dev/null
    echo "Processos finalizados."
    exit 0
}

# Captura Ctrl+C
trap cleanup SIGINT

# Mantém o script rodando
wait
