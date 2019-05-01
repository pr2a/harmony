export function levels() {
    return [
        {
          contents: [1, 0, 1, 1, 1, 1, 1, 1, 1],
          initialSelected: { x: 0, y: 0 }
        },
        {
          contents: [1, 1, 1, 1, 0, 1, 1, 0, 1],
          initialSelected: { x: 2, y: 2 }
        }
        // {
        //   contents: [1, -1, 0, 0, -1, 0, 0, 1, 0],
        //   initialSelected: { x: 1, y: 2 }
        // },
        // {
        //   contents: [-2, -4, -3, -1, -2, -4, -1, -1, -1],
        //   initialSelected: { x: 2, y: 0 }
        // },
        // {
        //   contents: [3, -1, -2, 3, -1, -3, 3, 0, -1],
        //   initialSelected: { x: 2, y: 0 }
        // }
      ];
}