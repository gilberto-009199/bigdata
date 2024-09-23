import org.apache.hadoop.io.{IntWritable, LongWritable, Text}
import org.apache.hadoop.mapreduce.Mapper

class ContadorDePalavrasMapper extends Mapper[LongWritable, Text, Text, IntWritable] {
  val one = new IntWritable(1)
  val word = new Text()

  override def map(key: LongWritable, value: Text, context: Mapper[LongWritable, Text, Text, IntWritable]#Context): Unit = {
    val line = value.toString
    val tokens = line.split("\\s+")
    tokens.foreach { token =>
      word.set(token)
      context.write(word, one)
    }
  }
}
