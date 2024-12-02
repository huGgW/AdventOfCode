use std::{
    error::{self},
    fs::File,
    io::{BufRead, BufReader},
};

fn main() {
    part1().expect("part1 failed");
    part2().expect("part2 failed");
}

fn part1() -> Result<(), Box<dyn error::Error>> {
    let reports = read_reports()?;

    let safe_cnt = reports
        .iter()
        .filter(|&r| is_safe(r))
        .collect::<Vec<_>>()
        .len();

    println!("Part 1: {}", safe_cnt);

    Ok(())
}

fn part2() -> Result<(), Box<dyn error::Error>> {
    let reports = read_reports()?;

    let safe_cnt = reports
        .iter()
        .filter(|&r| is_damper_safe(r))
        .collect::<Vec<_>>()
        .len();

    println!("Part 2: {}", safe_cnt);

    Ok(())
}

fn read_reports() -> Result<Vec<Vec<i32>>, Box<dyn error::Error>> {
    let file = File::open("input.txt")?;
    let reader = BufReader::new(file);

    reader
        .lines()
        .map(|line_res| -> Result<Vec<_>, Box<dyn error::Error>> {
            let line = line_res?;

            line.split_whitespace()
                .map(|s| s.parse::<i32>())
                .collect::<Result<Vec<_>, _>>()
                .map_err(|e| e.into())
        })
        .collect::<Result<Vec<_>, _>>()
}

fn is_safe(report: &[i32]) -> bool {
    assert!(report.len() >= 2, "report should have at least 2 elements");

    let is_safe_range = is_safe_range_fn(report[1] >= report[0]);

    (0..report.len() - 1).all(|i| is_safe_range(report[i], report[i + 1]))
}

fn is_damper_safe(report: &[i32]) -> bool {
    assert!(report.len() >= 3, "report should have at least 3 elements");

    let is_safe_range = is_safe_range_fn(report[1] >= report[0]);

    let mut error_start_idx: Option<usize> = Option::None;

    for i in 0..report.len() - 1 {
        if !is_safe_range(report[i], report[i + 1]) {
            error_start_idx = Option::Some(i);
            break;
        }
    }

    match error_start_idx {
        Some(error_start_idx) => if error_start_idx == 0 {
            vec![error_start_idx, error_start_idx + 1]
        } else {
            vec![error_start_idx, error_start_idx + 1, error_start_idx - 1]
        }
        .iter()
        .map(|&i| {
            let mut report = report.to_vec();
            let _ = report.remove(i);
            is_safe(&report)
        })
        .any(|b| b),

        None => true,
    }
}

fn is_safe_range_fn(is_ascending: bool) -> fn(i32, i32) -> bool {
    if is_ascending {
        |bef: i32, aft: i32| 1 <= aft - bef && aft - bef <= 3
    } else {
        |bef: i32, aft: i32| -3 <= aft - bef && aft - bef <= -1
    }
}
