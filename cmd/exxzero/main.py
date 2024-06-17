class DSU:
    def __init__(self, n):
        self.parent = [-1] * n
        self.rank = [1] * n

    def find(self, i):
        if self.parent[i] == -1:
            return i
        self.parent[i] = self.find(self.parent[i])
        return self.parent[i]

    def unite(self, x, y):
        s1 = self.find(x)
        s2 = self.find(y)
        
        if s1 != s2:
            if self.rank[s1] < self.rank[s2]:
                self.parent[s1] = s2
            elif self.rank[s1] > self.rank[s2]:
                self.parent[s2] = s1
            else:
                self.parent[s2] = s1
                self.rank[s1] += 1

class Graph:
    def __init__(self, V):
        self.edgelist = []
        self.V = V

    def add_edge(self, x, y, w):
        self.edgelist.append((w, x, y))

    def kruskals_mst(self):
        self.edgelist.sort()
        dsu = DSU(self.V)
        mst_cost = 0
        print("Following are the edges in the constructed MST")
        
        for w, x, y in self.edgelist:
            if dsu.find(x) != dsu.find(y):
                dsu.unite(x, y)
                mst_cost += w
                print(f"{x} -- {y} == {w}")
                
        print(f"Minimum Cost Spanning Tree: {mst_cost}")

def main():
    g = Graph(18)
    
    # Adding edges
    edges = [(1, 2, 1), (1, 3, 2), (1, 4, 4), (2, 5, 4), (2, 7, 2), (3, 6, 6), (3, 5, 7),
             (4, 6, 2), (4, 7, 3), (5, 8, 7), (5, 9, 5), (6, 8, 7), (6, 10, 3), (7, 9, 5),
             (7, 10, 4), (8, 11, 4), (9, 11, 1), (10, 11, 4)]
    for x, y, w in edges:
        g.add_edge(x, y, w)
    
    g.kruskals_mst()

if __name__ == "__main__":
    main()
