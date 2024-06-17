def tsp(adj_matrix, n, index_start):
    all_vertices = set(range(n))
    all_vertices.remove(index_start)
    current_vertex = index_start
    path = [current_vertex]
    weight = 0

    while all_vertices:
        next_vertex = min(all_vertices, key=lambda x: adj_matrix[current_vertex][x])
        all_vertices.remove(next_vertex)
        path.append(next_vertex)
        weight += adj_matrix[current_vertex][next_vertex]
        current_vertex = next_vertex

    path.append(index_start)
    weight += adj_matrix[current_vertex][index_start]
    return path, weight


def main():
    INF = float('inf')
    # matrix = [
    #     [ INF, 1, 5, 6, 5, 5, 1, 5 ],
    #     [ 1, INF, 7, 5, 6, 5, 5, 5 ],
    #     [ 5, 7, INF, 6, 1, 5, 6, 5 ],
    #     [ 6, 5, 6, INF, 7, 7, 7, 5 ],
    #     [ 5, 6, 1, 7, INF, 5, 6, 5 ],
    #     [ 5, 5, 5, 7, 5, INF, 5, 1 ],
    #     [ 1, 5, 6, 7, 6, 5, INF, 5 ],
    #     [ 5, 5, 5, 5, 5, 1, 5, INF ]
    # ]
        # matrix = [
    #     [inf, 1, 6, 7, 3, 5, 4, 6],
    #     [1, inf, 6, 5, 4, 2, 3, 5],
    #     [6, 6, inf, 1, 7, 6, 1, 5],
    #     [7, 5, 1, inf, 6, 7, 5, 1],
    #     [3, 4, 7, 6, inf, 5, 6, 5],
    #     [5, 2, 6, 7, 5, inf, 5, 4],
    #     [4, 3, 1, 5, 6, 5, inf, 7],
    #     [6, 5, 5, 1, 5, 4, 7, inf]
    # ]
    matrix = [
    [INF, 4, 6, 5, 5, 5, 4, 5],
    [4, INF, 7, 1, 6, 5, 1, 4],
    [6, 7, INF, 1, 1, 1, 5, 7],
    [5, 1, 1, INF, 6, 5, 1, 7],
    [5, 6, 1, 6, INF, 5, 5, 6],
    [5, 5, 1, 5, 5, INF, 5, 7],
    [4, 1, 5, 1, 5, 5, INF, 5],
    [5, 4, 7, 7, 6, 7, 5, INF]
]
    n = 8
    for i in range(n):
        path, weight = tsp(matrix, n, i)
        print(f'Path: {" → ".join([str(x + 1) for x in path])}, Weight: {weight}')


if __name__ == '__main__':
    main()
