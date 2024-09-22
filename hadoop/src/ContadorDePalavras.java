import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.Path;
import org.apache.hadoop.io.IntWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Job;
import org.apache.hadoop.mapreduce.lib.input.FileInputFormat;
import org.apache.hadoop.mapreduce.lib.output.FileOutputFormat;

public class ContadorDePalavras {
    public static void main(String[] args) throws Exception {
        // Cria a configuração do job
        Configuration conf = new Configuration();
        Job job = Job.getInstance(conf, "Contador de Palavras");

        // Define as classes para Mapper, Reducer e a classe de driver
        job.setJarByClass(ContadorDePalavras.class);
        job.setMapperClass(ContadorDePalavrasMapper.class);
        job.setReducerClass(ContadorDePalavrasReducer.class);

        // Define os tipos de saída do Mapper
        job.setMapOutputKeyClass(Text.class);
        job.setMapOutputValueClass(IntWritable.class);

        // Define os tipos de saída do Reducer
        job.setOutputKeyClass(Text.class);
        job.setOutputValueClass(IntWritable.class);

        // Define os caminhos de entrada e saída
        FileInputFormat.addInputPath(job, new Path(args[0]));
        FileOutputFormat.setOutputPath(job, new Path(args[1]));

        // Executa o job e aguarda a conclusão
        System.exit(job.waitForCompletion(true) ? 0 : 1);
    }
}
