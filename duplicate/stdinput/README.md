This piece of code prints each line that appears more than once inthe standard input preceded by it's count.
Each time the line is read from the input, the is used as a key into the map and the corresponding value is incremented.
It's not a problem if the map doesn't yet contain that key, the first a new line is seen the expression counts[line] on the right-hand side evaluates to the zero value for it's type.
Each iteration produces two results, a key and the value of the map element for that key