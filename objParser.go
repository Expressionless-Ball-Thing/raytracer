package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

// Parse Wavefront OBJ files to give a bunch of hittables.
// For now, it's just vertexes and faces.

type Parser struct {
	vertCoord []Vec3
	vertCount int
	list      Hit_List
}

func NewObj(filename string) *Hit_List {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	var parser Parser

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			// Last line?
			break
		} else {
			parser.parseVertexLine(line)
		}

	}

	return &parser.list

}

// Parse the line, see if it's only a vertex or a face.
func (parser *Parser) parseVertexLine(line string) {
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "v ") {
		temp := line[2:]
		list := strings.FieldsFunc(temp, unicode.IsSpace)
		result := make([]float64, len(list))
		for idx, str := range list {
			str = strings.TrimSpace(str)
			var err error
			if result[idx], err = strconv.ParseFloat(str, 64); err != nil {
				continue
			}
		}
		switch len(result) {
		case 3:
			parser.vertCoord = append(parser.vertCoord, *NewVec3(result[0], result[1], result[2]))
		case 4:
			w := result[3]
			parser.vertCoord = append(parser.vertCoord, *NewVec3(result[0]/w, result[1]/w, result[2]/w))
		default:
			fmt.Println("Malformed Vertex line:", line)
			os.Exit(1)
		}

	} else if strings.HasPrefix(line, "f ") {
		temp := line[2:]
		list := strings.FieldsFunc(temp, unicode.IsSpace)
		result := make([]string, len(list))
		for idx, str := range list {
			str = strings.TrimSpace(str)
			result[idx] = str
		}
		size := len(result)

		if size < 3 {
			fmt.Println("Malformed face line:", line)
			os.Exit(1)
		} else if size == 3 {
			triangle_incides := make([]int, 3)
			for idx, str := range result {
				index := strings.Split(str, "/")[0]
				if number, err := strconv.ParseInt(index, 0, 64); err != nil {
					fmt.Println(err)
					os.Exit(1)
				} else {
					triangle_incides[idx] = int(number - 1)
				}
			}
			vert1, vert2, vert3 := parser.vertCoord[triangle_incides[0]], parser.vertCoord[triangle_incides[1]], parser.vertCoord[triangle_incides[2]]

			parser.list.Add(NewTriangle(&vert1, vert2.Sub(&vert1), vert3.Sub(&vert1), NewLambert(*NewVec3(.12, .45, .15))))
		} else {
			// Decompose into a bunch of triangles

			triangle_incides := make([]int, size)
			for idx, str := range result {
				index := strings.Split(str, "/")[0]
				if number, err := strconv.ParseInt(index, 0, 64); err != nil {
					fmt.Println(err)
					os.Exit(1)
				} else {
					triangle_incides[idx] = int(number - 1)
				}
			}
			vert1 := parser.vertCoord[triangle_incides[0]]
			for index := 1; index < size-1; index++ {
				parser.list.Add(NewTriangle(&vert1, parser.vertCoord[triangle_incides[index]].Sub(&vert1), parser.vertCoord[triangle_incides[index+1]].Sub(&vert1), NewLambert(*NewVec3(.12, .45, .15))))
			}

		}
	}

	parser.vertCount++
}
