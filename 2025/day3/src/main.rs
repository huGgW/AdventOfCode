use std::{
    cmp::Reverse,
    collections::BinaryHeap,
    env,
    fs::File,
    io::{BufRead, BufReader},
};

fn main() {
    let input_file_name = env::args()
        .skip(1)
        .next()
        .expect("Please provide an input file name");
    let input_file = File::open(input_file_name).expect("Failed to open file");

    let batteries = parse_file(input_file);
    let two_largest_joltages: u32 = batteries.iter().map(Battery::two_largest_joltage).sum();

    let twelve_largest_joltages: u64 = batteries.iter().map(|b| b.n_largest_joltages(12)).sum();

    println!("Sum of two largest joltages: {}", two_largest_joltages);
    println!(
        "Sum of twelve largest joltages: {}",
        twelve_largest_joltages
    );
}

#[derive(Debug)]
pub struct Battery {
    digits: Vec<u8>,
}

impl Battery {
    pub fn parse_from(s: &str) -> Self {
        let digits = s
            .chars()
            .map(|c| c.to_string().parse().expect("invalid format"))
            .collect();

        Battery { digits }
    }

    pub fn two_largest_joltage(&self) -> u32 {
        assert!(self.digits.len() >= 2);

        let mut max_digit = None;
        let mut second_max_digit = None;

        let set_max_digit =
            |max_digit: &mut Option<u8>, second_max_digit: &mut Option<u8>, x: u8| {
                *max_digit = Some(x);
                *second_max_digit = None
            };

        let set_second_max_digit = |second_max_digit: &mut Option<u8>, x: u8| {
            *second_max_digit = Some(x);
        };

        for (i, &digit) in self.digits.iter().enumerate() {
            match (max_digit, second_max_digit) {
                (None, _) => set_max_digit(&mut max_digit, &mut second_max_digit, digit),
                (Some(x), _) if x < digit && i == self.digits.len() - 1 => {
                    set_second_max_digit(&mut second_max_digit, digit)
                }
                (Some(x), _) if x < digit => {
                    set_max_digit(&mut max_digit, &mut second_max_digit, digit)
                }
                (Some(x), None) if x >= digit => set_second_max_digit(&mut second_max_digit, digit),
                (Some(x), Some(y)) if x >= digit && y < digit => {
                    set_second_max_digit(&mut second_max_digit, digit)
                }
                _ => (),
            }
        }

        ((max_digit.unwrap()) as u32) * 10 + ((second_max_digit.unwrap()) as u32)
    }

    pub fn n_largest_joltages(&self, n: usize) -> u64 {
        assert!(n >= 1);
        assert!(self.digits.len() >= n);

        let mut heap = BinaryHeap::from_iter(
            self.digits
                .iter()
                .enumerate()
                .filter(|&(i, _)| i < self.digits.len() - n)
                .map(|(i, &d)| (d, Reverse(i))),
        );

        let mut result_vec = Vec::new();

        n_largest_joltages_recursive(&self.digits, n, 0, &mut heap, &mut result_vec)
    }
}

fn n_largest_joltages_recursive(
    digits: &[u8],
    n: usize,
    from: usize,
    heap: &mut BinaryHeap<(u8, Reverse<usize>)>,
    result_digits: &mut Vec<u8>,
) -> u64 {
    assert!(!heap.is_empty());
    assert!(digits.len() - from >= n as usize);

    if n == 0 {
        return calculate_number(result_digits);
    }

    if digits.len() - from == n as usize {
        digits[from..].iter().for_each(|&d| result_digits.push(d));
        return calculate_number(result_digits);
    }

    let before_idx = digits.len() - n as usize;

    // update heap
    heap.push((digits[before_idx - 1], Reverse(before_idx - 1)));
    while heap.peek().map(|&(_, i)| i.0 < from).unwrap_or(false) {
        heap.pop();
    }

    // get maximum value from left
    let (viable_max_digit, viable_max_idx) = heap.pop().unwrap();
    let viable_max_idx = viable_max_idx.0;

    if viable_max_digit < digits[before_idx] {
        digits[before_idx..]
            .iter()
            .for_each(|&d| result_digits.push(d));
        return calculate_number(result_digits);
    }

    result_digits.push(viable_max_digit);
    return n_largest_joltages_recursive(digits, n - 1, viable_max_idx + 1, heap, result_digits);
}

pub fn parse_file(file: File) -> Vec<Battery> {
    let buf_reader = BufReader::new(file);
    buf_reader
        .lines()
        .map(|r| r.expect("Failed to read line"))
        .filter(|l| !l.is_empty())
        .map(|l| Battery::parse_from(&l))
        .collect()
}

fn calculate_number(digits: &[u8]) -> u64 {
    let len = digits.len() as u32;
    digits.iter().enumerate().fold(0u64, |acc, (i, &d)| {
        10u64.pow(len - 1 - (i as u32)) * (d as u64) + acc
    })
}
