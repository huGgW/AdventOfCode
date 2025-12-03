use std::{error::Error, fs::File, io::Read};

fn main() {
    let args = std::env::args().collect::<Vec<String>>();
    let input_file = args.get(1).expect("Please provide an input file");

    let mut f = File::open(input_file).expect("Failed to open file");
    let mut buf = String::new();

    f.read_to_string(&mut buf).expect("Failed to read file");
    let id_ranges = parse(&buf).expect("Failed to parse input");

    let invalid_sum_1 = id_ranges
        .iter()
        .flat_map(|&(from, to)| invalid_ids_1(from, to))
        .sum::<u64>();

    println!("Sum of invalid ids of problem 1: {}", invalid_sum_1);

    let invalid_sum_2 = id_ranges
        .iter()
        .flat_map(|&(from, to)| invalid_ids_2(from, to))
        .sum::<u64>();

    println!("Sum of invalid ids of problem 2: {}", invalid_sum_2);
}

fn parse(input: &str) -> Result<Vec<(u64, u64)>, Box<dyn Error>> {
    input
        .trim()
        .split(",")
        .map(|range_str| -> Result<(u64, u64), Box<dyn Error>> {
            let sp = range_str.split("-").collect::<Vec<&str>>();

            if sp.len() != 2 {
                return Result::Err("Invalid range format".into());
            }

            let from = sp[0].parse::<u64>()?;
            let to = sp[1].parse::<u64>()?;

            Result::Ok((from, to))
        })
        .collect::<Result<Vec<(u64, u64)>, Box<dyn Error>>>()
}

fn invalid_ids_1(from: u64, to: u64) -> Vec<u64> {
    (from..=to).filter(|&x| is_repeated_by(x, 2)).collect()
}

fn invalid_ids_2(from: u64, to: u64) -> Vec<u64> {
    (from..=to).filter(|&x| is_repeated_any(x)).collect()
}

fn is_repeated_by(num: u64, count: u64) -> bool {
    digits_repeated_by(to_digits(num), count)
}

fn is_repeated_any(num: u64) -> bool {
    let digits = to_digits(num);
    let min_repeat_cnt = 2u64;
    let max_repeat_cnt = digits.len() as u64;

    (min_repeat_cnt..=max_repeat_cnt).any(|cnt| is_repeated_by(num, cnt))
}

fn digits_repeated_by(digits: Vec<u8>, count: u64) -> bool {
    if digits.len() % count as usize != 0 {
        return false;
    }

    let skip_len = digits.len() / count as usize;
    for i in 0..skip_len {
        let first_digit = digits[i];

        if (i..digits.len())
            .step_by(skip_len)
            .map(|j| digits[j])
            .any(|x| x != first_digit)
        {
            return false;
        }
    }

    true
}

fn to_digits(num: u64) -> Vec<u8> {
    let mut digit_vec = Vec::<u8>::new();

    let mut tmp_num = num;
    while tmp_num > 0 {
        let next_tmp_num = tmp_num / 10;
        let digit = tmp_num - next_tmp_num * 10;

        digit_vec.push(digit as u8);
        tmp_num = next_tmp_num;
    }

    digit_vec.reverse();
    digit_vec
}
