use std::{
    collections::{BinaryHeap, HashMap},
    fs::File,
    io::{BufRead, BufReader},
};

fn main() {
    part1();
    part2();
}

fn part1() {
    // Open input file
    let f = File::open("input.txt").expect("Cannot open input file");
    let reader = BufReader::new(f);

    // Create min heap
    let mut heap1 = BinaryHeap::new();
    let mut heap2 = BinaryHeap::new();

    // Collect two number to each heap
    reader.lines().for_each(|line_res| {
        let nums = line_res
            .expect("Cannot read line of file")
            .split_whitespace()
            .map(|s| s.parse::<i32>().expect("Cannot parse number"))
            .collect::<Vec<i32>>();

        heap1.push(std::cmp::Reverse(
            *(nums.get(0).expect("Should have first number")),
        ));

        heap2.push(std::cmp::Reverse(
            *(nums.get(1).expect("Should have second number")),
        ));
    });

    // Calculate distance
    let mut dist = 0;
    assert!(heap1.len() == heap2.len());
    while !heap1.is_empty() {
        let num1 = heap1.pop().unwrap().0;
        let num2 = heap2.pop().unwrap().0;

        dist += (num1 - num2).abs();
    }

    // Print out the answer
    println!("Part 1: {}", dist);
}

fn part2() {
    // Open input file
    let f = File::open("input.txt").expect("Cannot open input file");
    let reader = BufReader::new(f);

    // Create list for first numbers, counter for second numbers
    let mut nums = Vec::new();
    let mut counter = HashMap::new();

    // Collect from file
    reader.lines().for_each(|line_res| {
        let num_inputs = line_res
            .expect("Cannot read line of file")
            .split_whitespace()
            .map(|s| s.parse::<i32>().expect("Cannot parse number"))
            .collect::<Vec<i32>>();

        let &num1 = num_inputs.get(0).expect("Should have first number");
        let &num2 = num_inputs.get(1).expect("Should have second number");

        nums.push(num1);

        counter.insert(num2, *(counter.get(&num2).unwrap_or(&(0))) + 1);
    });

    // Calculate similarity score
    let similarity_score = nums
        .iter()
        .map(|num| *(counter.get(num).unwrap_or(&0)) * (*num))
        .sum::<i32>();

    // Print out the answer
    println!("Part 2: {}", similarity_score);
}
