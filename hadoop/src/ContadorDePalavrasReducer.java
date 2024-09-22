import org.apache.hadoop.io.IntWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Reducer;

import java.io.IOException;

public class ContadorDePalavrasReducer extends Reducer<Text, IntWritable, IntWritable, Text> {

    @Override
    protected void reduce(Text key, Iterable<IntWritable> values, Context context) throws IOException, InterruptedException {
        int sum = 0;

        // Soma todas as ocorrências da palavra
        for (IntWritable value : values) {
            sum += value.get();
        }

        // Emite a palavra com o total de ocorrências
        context.write(new IntWritable(sum), key);
    }
}
