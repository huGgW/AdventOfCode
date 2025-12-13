package org.example

import java.io.BufferedReader
import java.io.File
import java.io.InputStream
import java.io.InputStreamReader

fun main(args: Array<String>) {
    require(args.size == 1)
    val filepath = "../" + args.first()
    val file = File(filepath)

    file.inputStream().use { fs ->
        val (freshChecker, ids) = parseFrom(fs)

        val freshIds = ids.filter { freshChecker.isFresh(it) }
        println("fresh id count: ${freshIds.size}")

        val totalCount = freshChecker.totalFresh()
        println("total fresh id available count: $totalCount")
    }
}

fun parseFrom(`is`: InputStream): Pair<FreshChecker, List<Long>> {
    val br = BufferedReader(InputStreamReader(`is`))
    val fb = FreshChecker.builder()
    val ids = mutableListOf<Long>()

    // 0 : ranges, 1 : ids
    var parseStage = 0

    // Ranges
    while (true) {
        val line = br.readLine() ?: break
        if (line.isEmpty()) {
            assert(parseStage == 0)
            parseStage++
            continue
        }

        when (parseStage) {
            0 -> {
                val rangeVals = line.split("-").map { it.toLong() }
                require(rangeVals.size == 2)
                fb.addRange(rangeVals[0], rangeVals[1])
            }

            1 -> {
                line.toLong().let {
                    ids.add(it)
                }
            }
        }
    }

    return fb.build() to ids
}


