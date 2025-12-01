import file_streams/file_stream
import gleam/format
import gleam/int
import gleam/result
import gleam/string

pub fn main() -> Nil {
  let assert Ok(fs) = file_stream.open_read("input/part1.txt")
  // let assert Ok(fs) = file_stream.open_read("input/example1.txt")
  use <- defer(fn() {
    let assert Ok(_) = file_stream.close(fs)
    Nil
  })

  let locker = new_locker()

  let move_provider = fn() {
    file_stream.read_line(fs)
    |> result.map(string.trim)
    |> result.map(parse)
    |> result.map_error(fn(_) { Nil })
    |> result.flatten
  }

  let result_locker = move_loop(locker, move_provider)
  format.printf("password: ~b\n", result_locker.eventual_zero_cnt)
  format.printf("protocol password: ~b\n", result_locker.zero_cnt)
}

fn move_loop(locker: Locker, move_provider: fn() -> Result(Move, Nil)) -> Locker {
  case move_provider() {
    Ok(move_cmd) -> {
      let new_locker = move(locker, move_cmd)
      echo #(new_locker, move_cmd)
      move_loop(new_locker, move_provider)
    }
    Error(_) -> locker
  }
}

// Locker

type Locker {
  Locker(max_dial: Int, current: Int, eventual_zero_cnt: Int, zero_cnt: Int)
}

fn new_locker() -> Locker {
  Locker(max_dial: 99, current: 50, eventual_zero_cnt: 0, zero_cnt: 0)
}

fn move(locker: Locker, move: Move) -> Locker {
  case move {
    Left(steps) -> left_move(locker, steps)
    Right(steps) -> right_move(locker, steps)
  }
}

fn left_move(locker: Locker, steps: Int) -> Locker {
  case steps, locker.current {
    0, _ -> locker
    s, 0 if s > 0 -> {
      let mid_locker =
        Locker(
          locker.max_dial,
          locker.max_dial,
          locker.eventual_zero_cnt,
          locker.zero_cnt,
        )
      let remain_steps = s - 1

      left_move(mid_locker, remain_steps)
    }
    s, c if s == c ->
      Locker(
        locker.max_dial,
        0,
        locker.eventual_zero_cnt + 1,
        locker.zero_cnt + 1,
      )
    s, c if s < c ->
      Locker(locker.max_dial, c - s, locker.eventual_zero_cnt, locker.zero_cnt)
    s, c -> {
      // s > c
      let mid_locker =
        Locker(
          locker.max_dial,
          locker.max_dial,
          locker.eventual_zero_cnt,
          locker.zero_cnt + 1,
        )
      let remain_steps = s - c - 1

      left_move(mid_locker, remain_steps)
    }
  }
}

fn right_move(locker: Locker, steps: Int) -> Locker {
  case steps, locker.current {
    0, _ -> locker
    s, 0 if s > 0 -> {
      let mid_locker =
        Locker(locker.max_dial, 1, locker.eventual_zero_cnt, locker.zero_cnt)
      let remain_steps = s - 1

      right_move(mid_locker, remain_steps)
    }
    s, c if s + c == locker.max_dial + 1 ->
      Locker(
        locker.max_dial,
        0,
        locker.eventual_zero_cnt + 1,
        locker.zero_cnt + 1,
      )
    s, c if s + c < locker.max_dial + 1 -> {
      Locker(locker.max_dial, c + s, locker.eventual_zero_cnt, locker.zero_cnt)
    }
    s, c -> {
      // s + c > locker.max_dial + 1
      let mid_locker =
        Locker(
          locker.max_dial,
          1,
          locker.eventual_zero_cnt,
          locker.zero_cnt + 1,
        )
      let remain_steps = s - { locker.max_dial + 1 - c + 1 }
      right_move(mid_locker, remain_steps)
    }
  }
}

// Move

type Move {
  Left(steps: Int)
  Right(steps: Int)
}

fn parse(input: String) -> Result(Move, Nil) {
  case input {
    "L" <> rest -> int.parse(rest) |> result.map(Left)
    "R" <> rest -> int.parse(rest) |> result.map(Right)

    _ -> Error(Nil)
  }
}

// defer
pub fn defer(d: fn() -> Nil, f: fn() -> a) {
  let result = f()
  d()
  result
}
