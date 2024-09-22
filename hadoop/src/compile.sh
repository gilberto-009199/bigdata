#!/bin/sh

# Caminho onde o script está localizado (diretório atual do script)
PWDCOMPILE="$(cd "$(dirname "$0")" && pwd)"
# compila
../java-se-8u44-ri/bin/javac -cp "$PWDCOMPILE/hadoop-client-3.4.0.jar:$PWDCOMPILE/hadoop-common-3.4.0.jar:$PWDCOMPILE/hadoop-hdfs-3.4.0.jar:$PWDCOMPILE/hadoop-mapreduce-client-core-3.4.0.jar" *.java

# Verificar se a compilação foi bem-sucedida
if [ $? -eq 0 ]; then
    echo "Compilação bem-sucedida!"
    echo "Empacotando....."
    ../java-se-8u44-ri/bin/jar cf $PWDCOMPILE/contadorDePalavras.jar *.class
    if [ $? -eq 0 ]; then
        echo "Empacotamento bem sucedido: contadorDePalavras.jar"
    else
        echo "Erro no Empacotamento!"
    fi
else
    echo "Erro na compilação!"
fi