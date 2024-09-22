import org.apache.hadoop.io.IntWritable;
import org.apache.hadoop.io.LongWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Mapper;

import java.io.IOException;

public class ContadorDePalavrasMapper extends Mapper<LongWritable, Text, Text, IntWritable> {
    private final static IntWritable one = new IntWritable(1);
    private Text word = new Text();

    @Override
    protected void map(LongWritable key, Text value, Context context) throws IOException, InterruptedException {
        // Divide a linha de texto em palavras
        String[] words = value.toString().split("\\s+");

        // Emite cada palavra com o valor 1
        for (String str : words) {
            word.set(str);
            context.write(word, one);
        }
    }
}
