import java.io.*;
import java.util.*;

public class Solution {
    public static void main(String[] args) throws IOException {
        BufferedReader br = new BufferedReader(new InputStreamReader(System.in));
        int T = Integer.parseInt(br.readLine().trim());
        StringBuilder sb = new StringBuilder();

        while (T-- > 0) {
            StringTokenizer st = new StringTokenizer(br.readLine());
            int n = Integer.parseInt(st.nextToken());
            long d = Long.parseLong(st.nextToken());

            long[] b = new long[n];
            st = new StringTokenizer(br.readLine());
            for (int i = 0; i < n; i++) {
                b[i] = Long.parseLong(st.nextToken());
            }

            int insertions = 0;
            for (int i = 1; i < n; i++) {
                if (b[i] < b[i - 1] + d) {
                    insertions++;
                }
            }

            sb.append(n + insertions).append('\n');
            sb.append(b[0]);
            for (int i = 1; i < n; i++) {
                if (b[i] < b[i - 1] + d) {
                    sb.append(' ').append(1);
                }
                sb.append(' ').append(b[i]);
            }
            sb.append('\n');
        }

        System.out.print(sb);
    }
}
