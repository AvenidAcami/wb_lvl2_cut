package cmd

import (
	"fmt"
	"wb_lvl2_cut/internal/cut"

	"github.com/spf13/cobra"
)

var (
	fields    string
	delimiter string
	separated bool
)

var cutCmd = &cobra.Command{
	Use: "cut",
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := cut.GetResultChannel(cut.Options{
			Fields:    fields,
			Delimiter: delimiter,
			Separated: separated,
		})
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		for line := range result {
			fmt.Println(line)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cutCmd)
	cutCmd.Flags().StringVarP(&fields, "fields", "f", "", "Указание номеров полей (колонок), которые нужно вывести. Номера через запятую, можно диапазоны.")
	cutCmd.Flags().StringVarP(&delimiter, "delimiter", "d", "\t", "Использовать другой разделитель (символ). По умолчанию разделитель — табуляция ('\t')")
	cutCmd.Flags().BoolVarP(&separated, "separated", "s", false, "Только строки, содержащие разделитель. Если флаг указан, то строки без разделителя игнорируются (не выводятся).")
}
