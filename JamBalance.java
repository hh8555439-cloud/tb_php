import java.io.*;
import java.util.*;

public class JamBalance {
    public static void main(String[] args) throws IOException {
        BufferedReader br = new BufferedReader(new InputStreamReader(System.in));
        int T = Integer.parseInt(br.readLine().trim());
        StringBuilder sb = new StringBuilder();

        while (T-- > 0) {
            StringTokenizer st = new StringTokenizer(br.readLine());
            int n = Integer.parseInt(st.nextToken());
            int m = Integer.parseInt(st.nextToken());

            char[][] grid = new char[n][];
            int totalB = 0;
            for (int i = 0; i < n; i++) {
                grid[i] = br.readLine().trim().toCharArray();
                for (char c : grid[i]) {
                    if (c == 'B') totalB++;
                }
            }
            int totalS = n * m - totalB;
            int target = totalB - totalS;

            int offset = n * m;
            int size = 2 * offset + 1;
            final int INF = Integer.MAX_VALUE / 2;

            int[] dp = new int[size];
            Arrays.fill(dp, INF);
            dp[offset] = 0;

            for (int i = 0; i < n; i++) {
                Map<Integer, Integer> opts = new HashMap<>();
                int diff = 0;
                opts.put(0, 0);
                for (int k = 1; k <= m; k++) {
                    diff += (grid[i][k - 1] == 'B' ? 1 : -1);
                    opts.merge(diff, k, Math::min);
                }

                int[] newDp = new int[size];
                Arrays.fill(newDp, INF);

                for (int d = 0; d < size; d++) {
                    if (dp[d] >= INF) continue;
                    for (var e : opts.entrySet()) {
                        int nd = d + e.getKey();
                        if (nd >= 0 && nd < size) {
                            newDp[nd] = Math.min(newDp[nd], dp[d] + e.getValue());
                        }
                    }
                }

                dp = newDp;
            }

            int idx = offset + target;
            sb.append(idx >= 0 && idx < size && dp[idx] < INF ? dp[idx] : -1).append('\n');
        }

        System.out.print(sb);
    }
}
