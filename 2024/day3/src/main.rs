use std::{
    fs::File,
    io::{self, BufRead, BufReader, Read},
};

// State Machine
#[derive(Debug)]
enum State {
    M,
    U,
    L,
    D,
    O,
    N,
    T,
    Appostrophe,
    Comma,
    NumOne(u8),
    NumTwo(u8),
    MultOpenParen,
    DoOpenParen,
    DontOpenParen,
    MultCloseParen,
    DoCloseParen,
    DontCloseParen,
    Invalid,
}

fn state_change(state: &State, input: char) -> State {
    match (input, state) {
        // state changes for mult instructions
        ('m', _) => State::M,
        ('u', State::M) => State::U,
        ('l', State::U) => State::L,
        ('(', State::L) => State::MultOpenParen,
        (c, State::MultOpenParen) if c.is_numeric() => State::NumOne(0),
        (c, State::NumOne(l)) if c.is_numeric() && *l < 2 => State::NumOne(l + 1),
        (',', State::NumOne(l)) if *l <= 2 => State::Comma,
        (c, State::Comma) if c.is_numeric() => State::NumTwo(0),
        (c, State::NumTwo(l)) if c.is_numeric() && *l < 2 => State::NumTwo(l + 1),
        (')', State::NumTwo(l)) if *l <= 2 => State::MultCloseParen,

        // state changes for do, dont instructinos
        ('d', _) => State::D,
        ('o', State::D) => State::O,
        ('(', State::O) => State::DoOpenParen,
        (')', State::DoOpenParen) => State::DoCloseParen,
        ('n', State::O) => State::N,
        ('\'', State::N) => State::Appostrophe,
        ('t', State::Appostrophe) => State::T,
        ('(', State::T) => State::DontOpenParen,
        (')', State::DontOpenParen) => State::DontCloseParen,
        _ => State::Invalid,
    }
}

// Multiply Instruction
struct MulInstruct {
    x: i32,
    y: i32,
}
impl MulInstruct {
    fn new(s: String) -> MulInstruct {
        let num_strs: Vec<&str> = (s[4..s.len() - 1]).split(',').collect();
        let x = num_strs[0].parse::<i32>().unwrap();
        let y = num_strs[1].parse::<i32>().unwrap();
        MulInstruct { x, y }
    }

    fn calculate(&self) -> i32 {
        self.x * self.y
    }
}

fn main() {
    part_one().expect("Error in part one");
    part_two().expect("Error in part two");
}

fn part_one() -> Result<(), io::Error> {
    let reader = get_file_reader()?;
    let insts = reader
        .lines()
        .map(|s_res| -> Result<Vec<_>, _> {
            let s = s_res?;
            Result::Ok(get_mul_instructions(s, false))
        })
        .collect::<Result<Vec<_>, io::Error>>();

    let result = insts?
        .iter()
        .flatten()
        .map(|inst| inst.calculate())
        .sum::<i32>();

    println!("part 1: {}", result);
    Ok(())
}

fn part_two() -> Result<(), io::Error> {
    let mut reader = get_file_reader()?;
    let mut buffer = String::new();
    reader.read_to_string(&mut buffer)?;

    let insts = get_mul_instructions(buffer, true);
    let result = insts.iter().map(|inst| inst.calculate()).sum::<i32>();

    println!("part 2: {}", result);
    Ok(())
}

fn get_mul_instructions(s: String, use_toggle_inst: bool) -> Vec<MulInstruct> {
    let mut insts = Vec::new();
    let mut consider_mult = true;

    let mut inst_start_idx = Option::None;
    let mut state = State::Invalid;

    s.char_indices().for_each(|(i, c)| {
        state = state_change(&state, c);

        match state {
            State::M => inst_start_idx = Option::Some(i),
            State::MultCloseParen => {
                if consider_mult {
                    let inst = MulInstruct::new((s[inst_start_idx.unwrap()..=i]).to_string());
                    insts.push(inst);
                }
                inst_start_idx = Option::None;
            }

            State::D if use_toggle_inst => inst_start_idx = Option::Some(i),
            State::DoCloseParen if use_toggle_inst => {
                consider_mult = true;
                inst_start_idx = Option::None;
            }
            State::DontCloseParen if use_toggle_inst => {
                consider_mult = false;
                inst_start_idx = Option::None;
            }

            State::Invalid => {
                inst_start_idx = Option::None;
            }
            _ => {}
        }
    });

    insts
}

fn get_file_reader() -> Result<BufReader<File>, io::Error> {
    let file = File::open("input.txt")?;
    let reader = BufReader::new(file);
    Ok(reader)
}
