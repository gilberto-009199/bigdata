import org.apache.hadoop.conf.Configuration
import org.apache.hadoop.fs.Path
import org.apache.hadoop.io.IntWritable
import org.apache.hadoop.io.Text
import org.apache.hadoop.mapreduce.Job
import org.apache.hadoop.mapreduce.lib.input.FileInputFormat
import org.apache.hadoop.mapreduce.lib.output.FileOutputFormat

object ContadorDePalavras {
  def main(args: Array[String]): Unit = {
    // Cria a configuração do job
    val conf = new Configuration()
    val job = Job.getInstance(conf, "Contador de Palavras")

    // Define as classes para Mapper, Reducer e a classe de driver
    job.setJarByClass(getClass) // Use getClass para referência ao objeto
    job.setMapperClass(classOf[ContadorDePalavrasMapper])
    job.setReducerClass(classOf[ContadorDePalavrasReducer])

    // Define os tipos de saída do Mapper
    job.setMapOutputKeyClass(classOf[Text])
    job.setMapOutputValueClass(classOf[IntWritable])

    // Define os tipos de saída do Reducer
    job.setOutputKeyClass(classOf[Text])
    job.setOutputValueClass(classOf[IntWritable])

    // Define os caminhos de entrada e saída
    FileInputFormat.addInputPath(job, new Path(args(0)))
    FileOutputFormat.setOutputPath(job, new Path(args(1)))

    // Executa o job e aguarda a conclusão
    System.exit(if (job.waitForCompletion(true)) 0 else 1)
  }
}
