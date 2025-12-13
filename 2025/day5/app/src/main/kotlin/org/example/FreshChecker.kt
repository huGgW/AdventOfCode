package org.example

import java.util.ArrayDeque
import java.util.Deque
import kotlin.math.max
import kotlin.math.min

class FreshChecker private constructor(
    private var ranges: List<FreshRange>,
) {
    fun isFresh(id: Long) = ranges.any { it.isInRange(id) }

    fun totalFresh(): Long {
        compact()

        return ranges.sumOf { it.count() }
    }

    private fun compact() {
        val toBeCompacted: Deque<FreshRange> = ArrayDeque(ranges)
        val compactedRanges = mutableListOf<FreshRange>()

        while (toBeCompacted.isNotEmpty()) {
            val range = toBeCompacted.poll()!!

            var merged = false
            for ((i, existRange) in compactedRanges.withIndex()) {
                if (existRange.isMergable(range)) {
                    compactedRanges.removeAt(i)
                    toBeCompacted.offer(existRange merge range)
                    merged = true
                    break
                }
            }

            if (!merged) {
                compactedRanges.add(range)
            }
        }

        this.ranges = compactedRanges
    }

    companion object {
        fun builder() = FreshCheckerBuilder()

        class FreshCheckerBuilder {
            private val ranges = mutableListOf<FreshRange>()

            fun addRange(from: Long, to: Long): FreshCheckerBuilder {
                ranges.add(FreshRange(from, to))
                return this
            }

            fun build() = FreshChecker(ranges)
        }
    }
}

private data class FreshRange(
    val from: Long,
    val to: Long,
) {
    fun isInRange(id: Long) = id in from..to

    fun count() = to - from + 1

    fun isMergable(other: FreshRange) = !(this.to < other.from || other.to < this.from)

    infix fun merge(other: FreshRange): FreshRange {
        require(isMergable(other))

        return FreshRange(
            min(this.from, other.from),
            max(this.to, other.to),
        )
    }
}

