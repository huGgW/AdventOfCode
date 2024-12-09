use std::{
    fs::File,
    io::{self, BufRead, BufReader},
};

fn main() {
    let ans1 = part1().expect("Failed to read input");
    println!("Part 1: {}", ans1);

    let ans2 = part2().expect("Failed to read input");
    println!("Part 2: {}", ans2);
}

fn part1() -> Result<usize, io::Error> {
    let grid = read_input()?;
    let x_coords = get_coords_of(&grid, 'X');
    let total_xmas = x_coords
        .iter()
        .map(|x_coord| count_xmas_on(&grid, x_coord))
        .sum();

    Ok(total_xmas)
}

fn part2() -> Result<usize, io::Error> {
    let grid = read_input()?;
    let a_coords = get_coords_of(&grid, 'A');

    let count = a_coords.iter().filter(|a_coord| is_x_mas(&grid, a_coord)).count();

    Ok(count)
}

fn get_coords_of(grid: &[Vec<char>], letter: char) -> Vec<(usize, usize)> {
    grid
        .iter()
        .enumerate()
        .flat_map(|(i, col)| {
            col.iter()
                .enumerate()
                .filter_map(|(j, c)| {
                    if *c == letter {
                        Option::Some(j)
                    } else {
                        Option::None
                    }
                })
                .map(|j| (i, j))
                .collect::<Vec<_>>()
        })
        .collect::<Vec<_>>()
}

fn count_xmas_on(grid: &[Vec<char>], x_coord: &(usize, usize)) -> usize {
    let right_coords = (0..4)
        .map(|a| (x_coord.0 as isize, x_coord.1 as isize + a))
        .collect::<Vec<_>>();
    let down_coords = (0..4)
        .map(|a| (x_coord.0 as isize + a, x_coord.1 as isize))
        .collect::<Vec<_>>();
    let left_coords = (0..4)
        .map(|a| (x_coord.0 as isize, x_coord.1 as isize - a))
        .collect::<Vec<_>>();
    let up_coords = (0..4)
        .map(|a| (x_coord.0 as isize - a, x_coord.1 as isize))
        .collect::<Vec<_>>();
    let right_down_coords = (0..4)
        .map(|a| (x_coord.0 as isize + a, x_coord.1 as isize + a))
        .collect::<Vec<_>>();
    let right_up_coords = (0..4)
        .map(|a| (x_coord.0 as isize - a, x_coord.1 as isize + a))
        .collect::<Vec<_>>();
    let left_down_coords = (0..4)
        .map(|a| (x_coord.0 as isize + a, x_coord.1 as isize - a))
        .collect::<Vec<_>>();
    let left_up_coords = (0..4)
        .map(|a| (x_coord.0 as isize - a, x_coord.1 as isize - a))
        .collect::<Vec<_>>();

    [
        right_coords,
        down_coords,
        left_coords,
        up_coords,
        right_down_coords,
        right_up_coords,
        left_down_coords,
        left_up_coords,
    ]
    .iter()
    .map(|coords| {
        coords
            .iter()
            .map(|&(i, j)| grid.get(i as usize).and_then(|col| col.get(j as usize)))
            .collect::<Vec<_>>()
    })
    .filter(|vc| {
        vc == &vec![
            Option::Some(&'X'),
            Option::Some(&'M'),
            Option::Some(&'A'),
            Option::Some(&'S'),
        ]
    })
    .count()
}

fn is_x_mas(grid: &[Vec<char>], a_coord: &(usize, usize)) -> bool {
    let lt_rb_coords = [
        (a_coord.0 as isize - 1, a_coord.1 as isize - 1),
        (a_coord.0 as isize + 1, a_coord.1 as isize + 1),
    ];
    let lb_rt_coords = [
        (a_coord.0 as isize - 1, a_coord.1 as isize + 1),
        (a_coord.0 as isize + 1, a_coord.1 as isize - 1),
    ];

    [lt_rb_coords, lb_rt_coords].iter().all(|coords| {
        let [(i1, j1), (i2, j2)] = coords;
        let c1_opt = grid.get(*i1 as usize).and_then(|col| col.get(*j1 as usize));
        let c2_opt = grid.get(*i2 as usize).and_then(|col| col.get(*j2 as usize));

        (c1_opt, c2_opt) == (Option::Some(&'M'), Option::Some(&'S'))
        || (c1_opt, c2_opt) == (Option::Some(&'S'), Option::Some(&'M'))
    })
}

fn read_input() -> Result<Vec<Vec<char>>, io::Error> {
    let file = File::open("input.txt")?;
    let reader = BufReader::new(file);

    reader
        .lines()
        .map(|s_opt| {
            let s = s_opt?;
            Ok(s.chars().collect::<Vec<_>>())
        })
        .collect::<Result<Vec<_>, io::Error>>()
}
