#include <stdio.h>
#include <stdlib.h>
#include <string.h>

char* pi= "3.141592653589793238462643";

int pi_digit(int n) {
    int digit = pi[n] - '0';
    return digit;
}

typedef struct Node {
    int id;
    int pi_digit;
    struct Node **children;
    int num_children;
} Node;

Node *create_node(int id,int pi_digit) {
    Node *new_node = (Node *)malloc(sizeof(Node));
    if (new_node != NULL) {
        new_node->id = id;
        new_node->pi_digit = pi_digit;
        new_node->children = NULL;
        new_node->num_children = 0;
    }
    return new_node;
}

void build_tree(Node *root, int depth,int id) {

    if(depth+2>=strlen(pi)){
        return;
    }

    int num_children = pi_digit(depth+2);

    root->children = (Node **)malloc(num_children * sizeof(Node *));
    root->num_children = num_children;

    for (int i = 0; i < num_children; ++i) {
        root->children[i] = create_node(id+i+1,num_children);
    }
    
    if(num_children>0)
        build_tree(root->children[0], depth+1,id+num_children);
}

void  dfs_traverse(Node *root) {
    if (root == NULL) return;

    printf("%d ", root->id);

    for (int i = 0; i < root->num_children; ++i) {
        dfs_traverse(root->children[i]);
    }

}

void bfs_traverse(Node *root) {
    if (root == NULL) return;

    for(int i = 0; i < root->num_children; ++i) {
        printf("%d ", root->children[i]->id);
    }

    for(int i = 0; i < root->num_children; ++i) {
        bfs_traverse(root->children[i]);
    }
}


void free_tree(Node *root) {
    if (root == NULL) return;

    for (int i = 0; i < root->num_children; ++i) {
        free_tree(root->children[i]);
    }

    free(root->children);
    free(root);
}

int main() {
    Node *root = create_node(1,1);
    
    build_tree(root, 1,1);

    printf("BFS Traversal : 1 ");
    bfs_traverse(root);
    printf("\n");

    printf("DFS Traversal : ");
    dfs_traverse(root);
    printf("\n");
 
    free_tree(root);
    return 0;
}
