use core::error;
use std::{
    fs::File,
    io::{self, BufRead, BufReader, Lines, Read, Write},
};

struct Rule {
    before: i32,
    after: i32,
}
impl Rule {
    fn from(cond: String) -> Result<Rule, Box<dyn error::Error>> {
        let nums = cond
            .split('|')
            .map(|s| s.parse::<i32>())
            .collect::<Result<Vec<_>, _>>()?;

        if nums.len() != 2 {
            Err(Box::from("length should be 2"))
        } else {
            Ok(Rule {
                before: nums[0],
                after: nums[1],
            })
        }
    }

    fn check(&self, pages: &[i32]) -> bool {
        let bef_idx_opt = pages.iter().position(|&x| x == self.before);
        let aft_idx_opt = pages.iter().position(|&x| x == self.after);

        match (bef_idx_opt, aft_idx_opt) {
            (Some(bef_idx), Some(aft_idx)) => bef_idx < aft_idx,
            _ => true,
        }
    }

    fn swap(&self, pages: &mut [i32]) {
        if self.check(pages) {
            return;
        }

        let bef_idx = pages.iter().position(|&x| x == self.before).unwrap();
        let aft_idx = pages.iter().position(|&x| x == self.after).unwrap();
        pages.swap(bef_idx, aft_idx);
    }
}

fn main() {
    match part1() {
        Ok(v) => println!("part 1: {}", v),
        Err(e) => eprintln!("Error: {}", e),
    }

    match part2() {
        Ok(v) => println!("part 2: {}", v),
        Err(e) => eprintln!("Error: {}", e),
    }
}

fn part1() -> Result<i32, Box<dyn error::Error>> {
    let mut lines = lines()?;
    let rules = rules(&mut lines)?;
    let pages = pages(&mut lines)?;

    let valid_pages = pages
        .iter()
        .filter(|&v| rules.iter().all(|r| r.check(v)))
        .collect::<Vec<_>>();

    let middle_sum = valid_pages.iter().map(|&v| v[v.len() / 2]).sum();

    Ok(middle_sum)
}

fn part2() -> Result<i32, Box<dyn error::Error>> {
    let mut lines = lines()?;
    let rules = rules(&mut lines)?;
    let pages = pages(&mut lines)?;

    let mut invalid_pages_with_validity = pages
        .into_iter()
        .map(|v| {
            if rules.iter().all(|r| r.check(&v)) {
                (v, true)
            } else {
                (v, false)
            }
        })
        .filter(|(_, b)| !b)
        .collect::<Vec<_>>();

    while invalid_pages_with_validity.iter().any(|&(_, b)| !b) {
        invalid_pages_with_validity
            .iter_mut()
            .filter(|(_, b)| !b)
            .for_each(|(v, b)| {
                rules.iter().for_each(|r| r.swap(v));
                *b = rules.iter().all(|r| r.check(v));
            })
    }

    let middle_sum = invalid_pages_with_validity
        .iter()
        .map(|(v, _)| v[v.len() / 2])
        .sum();

    Ok(middle_sum)
}

fn lines() -> Result<Lines<BufReader<File>>, io::Error> {
    let file = File::open("input.txt")?;
    // let file = File::open("test.txt")?;
    let reader = BufReader::new(file);
    Ok(reader.lines())
}

fn rules<T: BufRead>(lines: &mut Lines<T>) -> Result<Vec<Rule>, Box<dyn error::Error>> {
    let mut rules = Vec::new();
    loop {
        let line = lines.next();
        match line {
            Some(Err(e)) => return Err(Box::new(e)),
            Some(Ok(s)) => match s {
                _ if s.is_empty() => break,
                _ => {
                    let rule = Rule::from(s)?;
                    rules.push(rule);
                }
            },
            _ => return Err(Box::from("seperate empty line should be given")),
        };
    }

    Ok(rules)
}

fn pages<T: BufRead>(lines: &mut Lines<T>) -> Result<Vec<Vec<i32>>, Box<dyn error::Error>> {
    let pages = lines
        .map(|s_res| -> Result<Vec<i32>, Box<dyn error::Error>> {
            let s = s_res?;
            s.split(',')
                .map(|ns| {
                    let n = ns.parse::<i32>()?;
                    Ok(n)
                })
                .collect::<Result<Vec<_>, _>>()
        })
        .collect::<Result<Vec<_>, _>>()?;

    Ok(pages)
}
