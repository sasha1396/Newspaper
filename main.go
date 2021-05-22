package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

func openfile(input *string) *os.File {
	fileopen, err := os.Open(*input)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//defer fileopen.Close()
	return fileopen
}

func createfile(output *string) *os.File {
	newfile, err := os.Create(*output) // создаем файл
	if err != nil {                    // если возникла ошибка
		fmt.Println("Unable to create file:", err)
		os.Exit(1) // выходим из программы
	}
	//defer newfile.Close() // закрываем файл
	fmt.Println("Create file successed")
	return newfile
}

//splitting on the strings with a given size
func splitstrings(fileopen *os.File, wdt *int, spr *int) []string {
	var str, str2 ([]string) // для форматированного текста
	var size, oldsize int    // для размера текущей строки
	reader := bufio.NewReader(fileopen)
	width := *wdt
	separ := *spr

	for {
		// читаем строки
		line, err := reader.ReadString('\n')
		// делим на слова
		words := strings.Fields(line)
		// преобразуем в тип byte
		buf := []byte(line)
		// проверка строки на превышение размера колонки
		if (utf8.RuneCount(buf)) > width {
			// пробегаемся по массиву чтобы проверить слова
			for i := 0; i < len(words); i++ {
				buf = []byte(words[i])
				// считаем размер всей строки из полученных слов
				size += utf8.RuneCount(buf)
				// проверяем каждое слово на превышение колонки
				if utf8.RuneCount(buf) > width {
					{
						fmt.Println("Size word bigger than width of column:", utf8.RuneCount(buf))
						return nil
					}
				}
				// проверка чтобы не превышало ширину столбца
				// если равна, то переходим к пробелам
				if size > width {
					if size != width {
						// Ставим размер данного слова, так как оно переносится на новую строку и будет там первым
						size = utf8.RuneCount(buf)
					} else {
						// сохраняем предыдущий размер
						oldsize = size
					}
					for k := (width - oldsize + separ); k > 1; k-- {
						// добавляем пробелы
						str = append(str, " ")
					}
					str2 = append(str2, (strings.Join(str, "")))
					str = nil //
				}
				// сохраняем предыдущий размер
				oldsize = size
				// добавляем слово в массив
				str = append(str, words[i])

				if (i + 1) < len(words) {
					// добавляем пробелы между словами
					str = append(str, " ")
					// + пробел
					size++
				}
			}

			// если закончилась считываемая строка, но колонка еще не заполнена
			for k := (width - oldsize + separ); k > 0; k-- {
				// добавляем пробелы
				str = append(str, " ")
			}
			str2 = append(str2, (strings.Join(str, "")))
			str = nil //
			size = 0
		} else {
			// пробегаемся по массиву чтобы проверить слова
			for i := 0; i < len(words); i++ {
				// считаем размер всей строки из полученных слов
				size += utf8.RuneCount([]byte(words[i]))
				// добавляем слово в массив
				str = append(str, words[i])
				if (i + 1) < len(words) {
					// добавляем пробелы между словами
					str = append(str, " ")
					// + пробел
					size++
				}
			}

			for k := (width - size + separ); k > 0; k-- {
				// добавляем пробелы
				str = append(str, " ")
			}
			str2 = append(str2, (strings.Join(str, "")))
			str = nil //
			size = 0
		}

		if err != nil {
			if err == io.EOF {

				break
			} else {
				fmt.Println(err)
				return nil
			}
		}
	}
	str2 = append(str2, (strings.Join(str, "")))
	return str2
}

//writing to file
func writetofile(newfile *os.File, count *int, col *int, str2 *[]string) *os.File {
	i := 0
	// для расчета строки которую нужно пихать
	pos := 0
	counter := *count
	columns := *col
	sizecolumn := (counter / columns)
	str := *str2

	// counter задает максимально кол-во строк, ведь столбцы не всегда могут быть одинакового размера, pos начинается с 0, поэтому count-1
	for pos != counter-1 {
		for curcolumn := 0; curcolumn < columns; curcolumn++ {
			pos = i + (curcolumn * sizecolumn)
			// проверка, если вдруг позиция не стала равно емкости массиву, а перешагнула ее
			if pos == i && i == sizecolumn {
				break
			}
			//если позиция выходит за массив
			if pos >= len(str) {
				break
			} else {
				newfile.WriteString(str[pos])
			}
		}
		// тоже самое, чтобы вывйти из 2-го for
		if pos == i && i == sizecolumn {
			break
		}
		i++
		// переход на новую строку
		newfile.WriteString("\n")
	}
	return newfile
}

func main() {
	//___________________________________________Command-Line Flags________________________________________//
	var inputfile string
	flag.StringVar(&inputfile, "infile", "война и мир.txt", "a name file for input file to edit it")
	var outputfile string
	flag.StringVar(&outputfile, "outfile", "output.txt", "a name file for create edited file")
	columnsPtr := flag.Int("columns", 20, "amount of columns in the edited text")
	widthPtr := flag.Int("width", 165, "size of columns in amount of character")
	flag.Parse()
	//_______________________________________________________________________________________________________

	//___________________________________________Print_options_of_args_______________________________________
	fmt.Println("input file:", inputfile)
	fmt.Println("output file:", outputfile)
	fmt.Println("columns:", *columnsPtr)
	fmt.Println("width:", *widthPtr)
	//_______________________________________________________________________________________________________

	//_______________________________________________Varialubles____________________________________________//
	var width int     // ширина столбца
	var columns int   // колличество столбцов
	var separ int = 3 // расстояние между колонками

	width = *widthPtr
	columns = *columnsPtr

	//create new file
	newfile := createfile(&outputfile)
	defer newfile.Close() // закрываем файл

	//open file
	fileopen := openfile(&inputfile)
	defer fileopen.Close()

	strings := splitstrings(fileopen, &width, &separ)
	counter := len(strings)
	fmt.Println("Количество строк", counter)

	//Cколько полных строк в массиве.
	oldcounter := counter

	//Чтобы во всех колонках было одинаковое
	//количество строк добавляя пустые.
	for counter%columns != 0 {
		counter++
	}

	fmt.Println("Дополнено строк: ", counter-oldcounter)
	// Чтобы текста хватило на заданные колонки
	// если нет - ошибка.
	if oldcounter/columns < 1 /*|| oldcounter%columns == 0*/ {
		fmt.Println("Задано слишком большое число колонок")
	} else {
		writetofile(newfile, &counter, &columns, &strings)
	}
}
