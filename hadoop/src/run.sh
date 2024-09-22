#!/bin/sh

echo "Running"
cd /
cat /mapred-site.xml > /opt/hadoop/etc/hadoop/mapred-site.xml
echo -e " + Setting mapred-site.xml && HADOOP_MAPRED_HOME"
export HADOOP_MAPRED_HOME=/opt/hadoop
echo -e " + Clear hdfs:/input"
hadoop fs -rm -r /input
echo -e " + Clear hdfs:/output"
hadoop fs -rm -r /output
echo -e " + Move input.txt to hdfs:/input"
hadoop fs -put /input.txt /input
hadoop fs -put /input.txt /input
echo -e " + Run hadoop jar contadorDePalavras.jar ContadorDePalavras /input /output"
hadoop jar contadorDePalavras.jar ContadorDePalavras /input /output
echo -e " + Result:"
hadoop fs -cat /output/part-r-00000