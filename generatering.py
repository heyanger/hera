import random
ls = []

n = 3
k = 5
r = 2

for i in range(n):
    for j in range(k):
        ls.append((random.randint(0, 180), i))

ls.sort() # sort by first

print(ls)

new_ls = []
for i in range(len(ls)):
    rand, node = ls[i]

    s = set([node])

    j = i
    while len(s) < r:
        j = (j + 1) % (n * k)

        s.add(ls[j][1])
    
    new_ls.append((rand, list(s)))

print(new_ls)