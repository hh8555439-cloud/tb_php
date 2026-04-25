import java.util.*;

public class AndClosure {
    public static void main(String[] args) {
        Scanner in = new Scanner(System.in);
        int T = in.nextInt();
        StringBuilder sb = new StringBuilder();

        while (T-- > 0) {
            int n = in.nextInt();
            int[] f = new int[1024];
            boolean[] g = new boolean[1024];
            Arrays.fill(f, 1023);

            for (int i = 0; i < n; i++) {
                int a = in.nextInt();
                f[a] &= a;
                g[a] = true;
            }

            for (int b = 0; b < 10; b++) {
                for (int v = 0; v < 1024; v++) {
                    if ((v & (1 << b)) == 0) {
                        f[v] &= f[v | (1 << b)];
                        g[v] |= g[v | (1 << b)];
                    }
                }
            }

            int count = 0;
            for (int v = 0; v < 1024; v++) {
                if (g[v] && f[v] == v) count++;
            }
            sb.append(count).append('\n');
        }

        System.out.print(sb);
    }
}
