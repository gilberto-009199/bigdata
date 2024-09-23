import org.apache.hadoop.io.{IntWritable, Text}
import org.apache.hadoop.mapreduce.Reducer
import scala.jdk.CollectionConverters._ 

class ContadorDePalavrasReducer extends Reducer[Text, IntWritable, IntWritable, Text] {
  override def reduce(key: Text, values: java.lang.Iterable[IntWritable], context: Reducer[Text, IntWritable, IntWritable, Text]#Context): Unit = {
    // Converte o java.util.Iterator em um Scala Iterator usando asScala
    val sum = values.iterator().asScala.foldLeft(0)(_ + _.get)
    context.write(new IntWritable(sum), key)
  }
}
