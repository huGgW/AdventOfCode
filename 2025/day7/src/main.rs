use anyhow::{anyhow, Context, Result};
use std::{
    collections::{HashMap, HashSet},
    env,
    fs::File,
    io::{BufRead, BufReader},
};

fn main() -> Result<()> {
    let filename = env::args()
        .nth(1)
        .context("filename should be given as args")?;

    let file = File::open(filename).context("file should be opened")?;
    let tachyon = parse_input(file).context("input should be parsed")?;

    let split_cnt = tachyon.split_cnt();
    println!("Number of splits: {}", split_cnt);

    let quantum_state_cnt = tachyon.quantum_state_cnt();
    println!("Number of quantum states: {}", quantum_state_cnt);

    Ok(())
}

#[derive(Debug)]
pub struct Tachyon {
    width: usize,
    height: usize,
    start: (usize, usize),
    splitters: HashMap<usize, HashSet<usize>>,
}

impl Tachyon {
    pub fn split_cnt(&self) -> usize {
        let mut splitted = 0usize;
        let mut current_beams = HashSet::from([self.start.1]);

        for i in self.start.0..self.height {
            let mut next_beams = HashSet::new();

            for &j in current_beams.iter() {
                if self.splitter_exists(i, j) {
                    splitted += 1;
                    next_beams.insert(j - 1);
                    next_beams.insert(j + 1);
                } else {
                    next_beams.insert(j);
                }
            }

            current_beams = next_beams;
        }

        splitted
    }

    pub fn quantum_state_cnt(&self) -> usize {
        let mut current_states = HashMap::from([(self.start.1, 1)]);

        for i in self.start.0..self.height {
            let mut next_states = HashMap::new();

            for (&j, &cnt) in current_states.iter() {
                let mut next_local_states = Vec::new();

                if self.splitter_exists(i, j) {
                    next_local_states.push((j - 1, cnt));
                    next_local_states.push((j + 1, cnt));
                } else {
                    next_local_states.push((j, cnt));
                }

                for (pos, cnt) in next_local_states {
                    *next_states.entry(pos).or_insert(0) += cnt;
                }
            }

            current_states = next_states;
        }

        current_states.values().sum()
    }

    fn splitter_exists(&self, i: usize, j: usize) -> bool {
        self.splitters
            .get(&i)
            .map(|s| s.contains(&j))
            .unwrap_or(false)
    }
}

pub fn parse_input(file: File) -> Result<Tachyon> {
    let reader = BufReader::new(file);

    let mut width_opt = None;
    let mut height = 0usize;
    let mut start_opt = None;
    let mut splitters = HashMap::new();

    for (i, line_result) in reader.lines().enumerate() {
        let line = line_result?;

        match width_opt {
            None => width_opt = Some(line.len()),
            Some(w) if w == line.len() => {}
            _ => return Err(anyhow!("Inconsistent line width")),
        };

        height += 1;

        let mut row_set = HashSet::new();
        for (j, c) in line.char_indices() {
            match c {
                'S' => {
                    if start_opt.is_none() {
                        start_opt = Some((i, j));
                    } else {
                        return Err(anyhow!("Multiple start positions found"));
                    }
                }
                '^' => {
                    row_set.insert(j);
                }
                '.' => {}
                _ => return Err(anyhow!("Invalid character in input")),
            }
        }
        splitters.insert(i, row_set);
    }

    Ok(Tachyon {
        width: width_opt.unwrap(),
        height,
        start: start_opt.unwrap(),
        splitters,
    })
}
