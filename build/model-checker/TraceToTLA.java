import java.io.FileInputStream;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;
import java.util.Scanner;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

import tlc2.value.ValueOutputStream;
import tlc2.value.impl.IntValue;
import tlc2.value.impl.RecordValue;
import tlc2.value.impl.StringValue;
import tlc2.value.impl.TupleValue;
import tlc2.value.impl.Value;
import util.UniqueString;

public class TraceToTLA {

	public static void main(String[] args) throws IOException {
		final Path in = Paths.get(args[0]);
		final Path out = Paths.get(args[1]);

		final int lineCount = (int) Files.lines(in).count();

		System.out.printf("Parsing %s lines in %s to %s.\n", lineCount, in, out);

		final Value[] v = new Value[lineCount];

		FileInputStream inputStream = null;
		Scanner sc = null;
		try {
			inputStream = new FileInputStream(in.toFile());
			sc = new Scanner(inputStream);
			int i = 0;
			while (sc.hasNextLine()) {
				final String line = sc.nextLine();
				v[i++] = getValue(line);
			}
			if (sc.ioException() != null) {
				throw sc.ioException();
			}
		} finally {
			if (inputStream != null) {
				inputStream.close();
			}
			if (sc != null) {
				sc.close();
			}
		}

		final ValueOutputStream vos = new ValueOutputStream(out.toFile());
		// Do not normalize TupleValue because normalization depends on the actual UniqueString#internTable.
		new TupleValue(v).write(vos);
		vos.close();

		System.out.printf("Successfully parsed %s to %s.\n", in, out);
	}

	private static final Pattern f2 = Pattern.compile("<<(.*?)>>,<<(.*?)>>,<<(.*?)>>");

	private static final TupleValue getValue(String line) {
		final Value[] values = new Value[6];

		// Skip first tuple marker.
		line = line.substring(2);

		// field 0:
		values[0] = IntValue.gen(Integer.parseInt(line.substring(0, 1)));
		line = line.substring(3);

		// field 1:
		int indexOf = line.indexOf('"');
		values[1] = new StringValue(line.substring(0, indexOf));
		line = line.substring(indexOf + 2);

		// field 2: find end of first triple.
		indexOf = line.indexOf(">>>>");
		final String triple = line.substring(2, indexOf + 2);
		line = line.substring(indexOf + 5);

		final Matcher matcher = f2.matcher(triple);
		matcher.find();

		final Value[] triplet = new Value[3];
		triplet[0] = tupleToRecordsWrap(matcher.group(1));
		triplet[1] = tupleToRecordsWrap(matcher.group(2));
		triplet[2] = tupleToRecordsWrap(matcher.group(3));
		values[2] = new TupleValue(triplet);

		// field 3: <<"Leader","Follower","Follower">>,
		indexOf = line.indexOf(">>");
		final String substring = line.substring(2, indexOf);
		line = line.substring(indexOf + 3);
		final String[] split2 = substring.split(",");
		final Value[] stngs = new Value[3];
		for (int i = 0; i < split2.length; i++) {
			stngs[i] = new StringValue(split2[i].replace("\"", ""));
		}
		values[3] = new TupleValue(stngs);

		// field 4:
		indexOf = line.indexOf(">>");
		values[4] = tupleToRecords(line.substring(0, indexOf + 2));
		line = line.substring(indexOf + 3);

		// field 5:
		indexOf = line.lastIndexOf('"');
		values[5] = new StringValue(line.substring(1, indexOf));

		return new TupleValue(values);
	}

	private static final Pattern pattern = Pattern.compile("(\\[.*?\\])");

	private static Value tupleToRecordsWrap(String tuple) {
		return tupleToRecords("<<" + tuple + ">>");
	}

	private static Value tupleToRecords(String tuple) {
		if ("<<>>".equals(tuple)) {
			return TupleValue.EmptyTuple;
		}

		final List<Value> v = new ArrayList<>();

		final Matcher matcher = pattern.matcher(tuple.substring(2, tuple.length() - 2));
		while (matcher.find()) {
			final String group = matcher.group();
			v.add(toRecord(group));
		}
		return new TupleValue(v.toArray(new Value[v.size()]));
	}

	private static Value toRecord(String rcd) {
		rcd = rcd.substring(1, rcd.length() - 1);
		final String[] split = rcd.split(",");

		final UniqueString[] keys = new UniqueString[split.length];
		final Value[] vals = new Value[split.length];

		for (int i = 0; i < split.length; i++) {
			final String[] kv = split[i].trim().split(" \\|-> ");
			keys[i] = UniqueString.uniqueStringOf(kv[0]);
			vals[i] = IntValue.gen(Integer.parseInt(kv[1]));
		}
		return new RecordValue(keys, vals, false);
	}
}