// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'local_measurement.dart';

// **************************************************************************
// IsarCollectionGenerator
// **************************************************************************

// coverage:ignore-file
// ignore_for_file: duplicate_ignore, non_constant_identifier_names, constant_identifier_names, invalid_use_of_protected_member, unnecessary_cast, prefer_const_constructors, lines_longer_than_80_chars, require_trailing_commas, inference_failure_on_function_invocation, unnecessary_parenthesis, unnecessary_raw_strings, unnecessary_null_checks, join_return_with_assignment, prefer_final_locals, avoid_js_rounded_ints, avoid_positional_boolean_parameters, always_specify_types

extension GetLocalMeasurementCollection on Isar {
  IsarCollection<LocalMeasurement> get localMeasurements => this.collection();
}

const LocalMeasurementSchema = CollectionSchema(
  name: r'LocalMeasurement',
  id: -1357802536738280192,
  properties: {
    r'clientLocalId': PropertySchema(
      id: 0,
      name: r'clientLocalId',
      type: IsarType.string,
    ),
    r'deviceMac': PropertySchema(
      id: 1,
      name: r'deviceMac',
      type: IsarType.string,
    ),
    r'diffSignal': PropertySchema(
      id: 2,
      name: r'diffSignal',
      type: IsarType.doubleList,
    ),
    r'fingerprint': PropertySchema(
      id: 3,
      name: r'fingerprint',
      type: IsarType.doubleList,
    ),
    r'healthScore': PropertySchema(
      id: 4,
      name: r'healthScore',
      type: IsarType.long,
    ),
    r'isDeleted': PropertySchema(
      id: 5,
      name: r'isDeleted',
      type: IsarType.bool,
    ),
    r'isSynced': PropertySchema(
      id: 6,
      name: r'isSynced',
      type: IsarType.bool,
    ),
    r'measuredAt': PropertySchema(
      id: 7,
      name: r'measuredAt',
      type: IsarType.dateTime,
    ),
    r'riskLabel': PropertySchema(
      id: 8,
      name: r'riskLabel',
      type: IsarType.string,
    ),
    r'serverId': PropertySchema(
      id: 9,
      name: r'serverId',
      type: IsarType.string,
    ),
    r'syncedAt': PropertySchema(
      id: 10,
      name: r'syncedAt',
      type: IsarType.dateTime,
    ),
    r'updatedAt': PropertySchema(
      id: 11,
      name: r'updatedAt',
      type: IsarType.dateTime,
    )
  },
  estimateSize: _localMeasurementEstimateSize,
  serialize: _localMeasurementSerialize,
  deserialize: _localMeasurementDeserialize,
  deserializeProp: _localMeasurementDeserializeProp,
  idName: r'id',
  indexes: {
    r'clientLocalId': IndexSchema(
      id: 7868725324801303552,
      name: r'clientLocalId',
      unique: true,
      replace: false,
      properties: [
        IndexPropertySchema(
          name: r'clientLocalId',
          type: IndexType.hash,
          caseSensitive: true,
        )
      ],
    ),
    r'isSynced': IndexSchema(
      id: -39763503327887512,
      name: r'isSynced',
      unique: false,
      replace: false,
      properties: [
        IndexPropertySchema(
          name: r'isSynced',
          type: IndexType.value,
          caseSensitive: false,
        )
      ],
    )
  },
  links: {},
  embeddedSchemas: {},
  getId: _localMeasurementGetId,
  getLinks: _localMeasurementGetLinks,
  attach: _localMeasurementAttach,
  version: '3.1.0+1',
);

int _localMeasurementEstimateSize(
  LocalMeasurement object,
  List<int> offsets,
  Map<Type, List<int>> allOffsets,
) {
  var bytesCount = offsets.last;
  bytesCount += 3 + object.clientLocalId.length * 3;
  bytesCount += 3 + object.deviceMac.length * 3;
  bytesCount += 3 + object.diffSignal.length * 8;
  bytesCount += 3 + object.fingerprint.length * 8;
  bytesCount += 3 + object.riskLabel.length * 3;
  {
    final value = object.serverId;
    if (value != null) {
      bytesCount += 3 + value.length * 3;
    }
  }
  return bytesCount;
}

void _localMeasurementSerialize(
  LocalMeasurement object,
  IsarWriter writer,
  List<int> offsets,
  Map<Type, List<int>> allOffsets,
) {
  writer.writeString(offsets[0], object.clientLocalId);
  writer.writeString(offsets[1], object.deviceMac);
  writer.writeDoubleList(offsets[2], object.diffSignal);
  writer.writeDoubleList(offsets[3], object.fingerprint);
  writer.writeLong(offsets[4], object.healthScore);
  writer.writeBool(offsets[5], object.isDeleted);
  writer.writeBool(offsets[6], object.isSynced);
  writer.writeDateTime(offsets[7], object.measuredAt);
  writer.writeString(offsets[8], object.riskLabel);
  writer.writeString(offsets[9], object.serverId);
  writer.writeDateTime(offsets[10], object.syncedAt);
  writer.writeDateTime(offsets[11], object.updatedAt);
}

LocalMeasurement _localMeasurementDeserialize(
  Id id,
  IsarReader reader,
  List<int> offsets,
  Map<Type, List<int>> allOffsets,
) {
  final object = LocalMeasurement();
  object.clientLocalId = reader.readString(offsets[0]);
  object.deviceMac = reader.readString(offsets[1]);
  object.diffSignal = reader.readDoubleList(offsets[2]) ?? [];
  object.fingerprint = reader.readDoubleList(offsets[3]) ?? [];
  object.healthScore = reader.readLong(offsets[4]);
  object.id = id;
  object.isDeleted = reader.readBool(offsets[5]);
  object.isSynced = reader.readBool(offsets[6]);
  object.measuredAt = reader.readDateTime(offsets[7]);
  object.riskLabel = reader.readString(offsets[8]);
  object.serverId = reader.readStringOrNull(offsets[9]);
  object.syncedAt = reader.readDateTimeOrNull(offsets[10]);
  object.updatedAt = reader.readDateTime(offsets[11]);
  return object;
}

P _localMeasurementDeserializeProp<P>(
  IsarReader reader,
  int propertyId,
  int offset,
  Map<Type, List<int>> allOffsets,
) {
  switch (propertyId) {
    case 0:
      return (reader.readString(offset)) as P;
    case 1:
      return (reader.readString(offset)) as P;
    case 2:
      return (reader.readDoubleList(offset) ?? []) as P;
    case 3:
      return (reader.readDoubleList(offset) ?? []) as P;
    case 4:
      return (reader.readLong(offset)) as P;
    case 5:
      return (reader.readBool(offset)) as P;
    case 6:
      return (reader.readBool(offset)) as P;
    case 7:
      return (reader.readDateTime(offset)) as P;
    case 8:
      return (reader.readString(offset)) as P;
    case 9:
      return (reader.readStringOrNull(offset)) as P;
    case 10:
      return (reader.readDateTimeOrNull(offset)) as P;
    case 11:
      return (reader.readDateTime(offset)) as P;
    default:
      throw IsarError('Unknown property with id $propertyId');
  }
}

Id _localMeasurementGetId(LocalMeasurement object) {
  return object.id;
}

List<IsarLinkBase<dynamic>> _localMeasurementGetLinks(LocalMeasurement object) {
  return [];
}

void _localMeasurementAttach(
    IsarCollection<dynamic> col, Id id, LocalMeasurement object) {
  object.id = id;
}

extension LocalMeasurementByIndex on IsarCollection<LocalMeasurement> {
  Future<LocalMeasurement?> getByClientLocalId(String clientLocalId) {
    return getByIndex(r'clientLocalId', [clientLocalId]);
  }

  LocalMeasurement? getByClientLocalIdSync(String clientLocalId) {
    return getByIndexSync(r'clientLocalId', [clientLocalId]);
  }

  Future<bool> deleteByClientLocalId(String clientLocalId) {
    return deleteByIndex(r'clientLocalId', [clientLocalId]);
  }

  bool deleteByClientLocalIdSync(String clientLocalId) {
    return deleteByIndexSync(r'clientLocalId', [clientLocalId]);
  }

  Future<List<LocalMeasurement?>> getAllByClientLocalId(
      List<String> clientLocalIdValues) {
    final values = clientLocalIdValues.map((e) => [e]).toList();
    return getAllByIndex(r'clientLocalId', values);
  }

  List<LocalMeasurement?> getAllByClientLocalIdSync(
      List<String> clientLocalIdValues) {
    final values = clientLocalIdValues.map((e) => [e]).toList();
    return getAllByIndexSync(r'clientLocalId', values);
  }

  Future<int> deleteAllByClientLocalId(List<String> clientLocalIdValues) {
    final values = clientLocalIdValues.map((e) => [e]).toList();
    return deleteAllByIndex(r'clientLocalId', values);
  }

  int deleteAllByClientLocalIdSync(List<String> clientLocalIdValues) {
    final values = clientLocalIdValues.map((e) => [e]).toList();
    return deleteAllByIndexSync(r'clientLocalId', values);
  }

  Future<Id> putByClientLocalId(LocalMeasurement object) {
    return putByIndex(r'clientLocalId', object);
  }

  Id putByClientLocalIdSync(LocalMeasurement object, {bool saveLinks = true}) {
    return putByIndexSync(r'clientLocalId', object, saveLinks: saveLinks);
  }

  Future<List<Id>> putAllByClientLocalId(List<LocalMeasurement> objects) {
    return putAllByIndex(r'clientLocalId', objects);
  }

  List<Id> putAllByClientLocalIdSync(List<LocalMeasurement> objects,
      {bool saveLinks = true}) {
    return putAllByIndexSync(r'clientLocalId', objects, saveLinks: saveLinks);
  }
}

extension LocalMeasurementQueryWhereSort
    on QueryBuilder<LocalMeasurement, LocalMeasurement, QWhere> {
  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterWhere> anyId() {
    return QueryBuilder.apply(this, (query) {
      return query.addWhereClause(const IdWhereClause.any());
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterWhere> anyIsSynced() {
    return QueryBuilder.apply(this, (query) {
      return query.addWhereClause(
        const IndexWhereClause.any(indexName: r'isSynced'),
      );
    });
  }
}

extension LocalMeasurementQueryWhere
    on QueryBuilder<LocalMeasurement, LocalMeasurement, QWhereClause> {
  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterWhereClause> idEqualTo(
      Id id) {
    return QueryBuilder.apply(this, (query) {
      return query.addWhereClause(IdWhereClause.between(
        lower: id,
        upper: id,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterWhereClause>
      idNotEqualTo(Id id) {
    return QueryBuilder.apply(this, (query) {
      if (query.whereSort == Sort.asc) {
        return query
            .addWhereClause(
              IdWhereClause.lessThan(upper: id, includeUpper: false),
            )
            .addWhereClause(
              IdWhereClause.greaterThan(lower: id, includeLower: false),
            );
      } else {
        return query
            .addWhereClause(
              IdWhereClause.greaterThan(lower: id, includeLower: false),
            )
            .addWhereClause(
              IdWhereClause.lessThan(upper: id, includeUpper: false),
            );
      }
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterWhereClause>
      idGreaterThan(Id id, {bool include = false}) {
    return QueryBuilder.apply(this, (query) {
      return query.addWhereClause(
        IdWhereClause.greaterThan(lower: id, includeLower: include),
      );
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterWhereClause>
      idLessThan(Id id, {bool include = false}) {
    return QueryBuilder.apply(this, (query) {
      return query.addWhereClause(
        IdWhereClause.lessThan(upper: id, includeUpper: include),
      );
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterWhereClause> idBetween(
    Id lowerId,
    Id upperId, {
    bool includeLower = true,
    bool includeUpper = true,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addWhereClause(IdWhereClause.between(
        lower: lowerId,
        includeLower: includeLower,
        upper: upperId,
        includeUpper: includeUpper,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterWhereClause>
      clientLocalIdEqualTo(String clientLocalId) {
    return QueryBuilder.apply(this, (query) {
      return query.addWhereClause(IndexWhereClause.equalTo(
        indexName: r'clientLocalId',
        value: [clientLocalId],
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterWhereClause>
      clientLocalIdNotEqualTo(String clientLocalId) {
    return QueryBuilder.apply(this, (query) {
      if (query.whereSort == Sort.asc) {
        return query
            .addWhereClause(IndexWhereClause.between(
              indexName: r'clientLocalId',
              lower: [],
              upper: [clientLocalId],
              includeUpper: false,
            ))
            .addWhereClause(IndexWhereClause.between(
              indexName: r'clientLocalId',
              lower: [clientLocalId],
              includeLower: false,
              upper: [],
            ));
      } else {
        return query
            .addWhereClause(IndexWhereClause.between(
              indexName: r'clientLocalId',
              lower: [clientLocalId],
              includeLower: false,
              upper: [],
            ))
            .addWhereClause(IndexWhereClause.between(
              indexName: r'clientLocalId',
              lower: [],
              upper: [clientLocalId],
              includeUpper: false,
            ));
      }
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterWhereClause>
      isSyncedEqualTo(bool isSynced) {
    return QueryBuilder.apply(this, (query) {
      return query.addWhereClause(IndexWhereClause.equalTo(
        indexName: r'isSynced',
        value: [isSynced],
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterWhereClause>
      isSyncedNotEqualTo(bool isSynced) {
    return QueryBuilder.apply(this, (query) {
      if (query.whereSort == Sort.asc) {
        return query
            .addWhereClause(IndexWhereClause.between(
              indexName: r'isSynced',
              lower: [],
              upper: [isSynced],
              includeUpper: false,
            ))
            .addWhereClause(IndexWhereClause.between(
              indexName: r'isSynced',
              lower: [isSynced],
              includeLower: false,
              upper: [],
            ));
      } else {
        return query
            .addWhereClause(IndexWhereClause.between(
              indexName: r'isSynced',
              lower: [isSynced],
              includeLower: false,
              upper: [],
            ))
            .addWhereClause(IndexWhereClause.between(
              indexName: r'isSynced',
              lower: [],
              upper: [isSynced],
              includeUpper: false,
            ));
      }
    });
  }
}

extension LocalMeasurementQueryFilter
    on QueryBuilder<LocalMeasurement, LocalMeasurement, QFilterCondition> {
  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      clientLocalIdEqualTo(
    String value, {
    bool caseSensitive = true,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.equalTo(
        property: r'clientLocalId',
        value: value,
        caseSensitive: caseSensitive,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      clientLocalIdGreaterThan(
    String value, {
    bool include = false,
    bool caseSensitive = true,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.greaterThan(
        include: include,
        property: r'clientLocalId',
        value: value,
        caseSensitive: caseSensitive,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      clientLocalIdLessThan(
    String value, {
    bool include = false,
    bool caseSensitive = true,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.lessThan(
        include: include,
        property: r'clientLocalId',
        value: value,
        caseSensitive: caseSensitive,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      clientLocalIdBetween(
    String lower,
    String upper, {
    bool includeLower = true,
    bool includeUpper = true,
    bool caseSensitive = true,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.between(
        property: r'clientLocalId',
        lower: lower,
        includeLower: includeLower,
        upper: upper,
        includeUpper: includeUpper,
        caseSensitive: caseSensitive,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      clientLocalIdStartsWith(
    String value, {
    bool caseSensitive = true,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.startsWith(
        property: r'clientLocalId',
        value: value,
        caseSensitive: caseSensitive,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      clientLocalIdEndsWith(
    String value, {
    bool caseSensitive = true,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.endsWith(
        property: r'clientLocalId',
        value: value,
        caseSensitive: caseSensitive,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      clientLocalIdContains(String value, {bool caseSensitive = true}) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.contains(
        property: r'clientLocalId',
        value: value,
        caseSensitive: caseSensitive,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      clientLocalIdMatches(String pattern, {bool caseSensitive = true}) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.matches(
        property: r'clientLocalId',
        wildcard: pattern,
        caseSensitive: caseSensitive,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      clientLocalIdIsEmpty() {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.equalTo(
        property: r'clientLocalId',
        value: '',
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      clientLocalIdIsNotEmpty() {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.greaterThan(
        property: r'clientLocalId',
        value: '',
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      deviceMacEqualTo(
    String value, {
    bool caseSensitive = true,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.equalTo(
        property: r'deviceMac',
        value: value,
        caseSensitive: caseSensitive,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      deviceMacGreaterThan(
    String value, {
    bool include = false,
    bool caseSensitive = true,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.greaterThan(
        include: include,
        property: r'deviceMac',
        value: value,
        caseSensitive: caseSensitive,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      deviceMacLessThan(
    String value, {
    bool include = false,
    bool caseSensitive = true,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.lessThan(
        include: include,
        property: r'deviceMac',
        value: value,
        caseSensitive: caseSensitive,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      deviceMacBetween(
    String lower,
    String upper, {
    bool includeLower = true,
    bool includeUpper = true,
    bool caseSensitive = true,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.between(
        property: r'deviceMac',
        lower: lower,
        includeLower: includeLower,
        upper: upper,
        includeUpper: includeUpper,
        caseSensitive: caseSensitive,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      deviceMacStartsWith(
    String value, {
    bool caseSensitive = true,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.startsWith(
        property: r'deviceMac',
        value: value,
        caseSensitive: caseSensitive,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      deviceMacEndsWith(
    String value, {
    bool caseSensitive = true,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.endsWith(
        property: r'deviceMac',
        value: value,
        caseSensitive: caseSensitive,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      deviceMacContains(String value, {bool caseSensitive = true}) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.contains(
        property: r'deviceMac',
        value: value,
        caseSensitive: caseSensitive,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      deviceMacMatches(String pattern, {bool caseSensitive = true}) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.matches(
        property: r'deviceMac',
        wildcard: pattern,
        caseSensitive: caseSensitive,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      deviceMacIsEmpty() {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.equalTo(
        property: r'deviceMac',
        value: '',
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      deviceMacIsNotEmpty() {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.greaterThan(
        property: r'deviceMac',
        value: '',
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      diffSignalElementEqualTo(
    double value, {
    double epsilon = Query.epsilon,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.equalTo(
        property: r'diffSignal',
        value: value,
        epsilon: epsilon,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      diffSignalElementGreaterThan(
    double value, {
    bool include = false,
    double epsilon = Query.epsilon,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.greaterThan(
        include: include,
        property: r'diffSignal',
        value: value,
        epsilon: epsilon,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      diffSignalElementLessThan(
    double value, {
    bool include = false,
    double epsilon = Query.epsilon,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.lessThan(
        include: include,
        property: r'diffSignal',
        value: value,
        epsilon: epsilon,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      diffSignalElementBetween(
    double lower,
    double upper, {
    bool includeLower = true,
    bool includeUpper = true,
    double epsilon = Query.epsilon,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.between(
        property: r'diffSignal',
        lower: lower,
        includeLower: includeLower,
        upper: upper,
        includeUpper: includeUpper,
        epsilon: epsilon,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      diffSignalLengthEqualTo(int length) {
    return QueryBuilder.apply(this, (query) {
      return query.listLength(
        r'diffSignal',
        length,
        true,
        length,
        true,
      );
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      diffSignalIsEmpty() {
    return QueryBuilder.apply(this, (query) {
      return query.listLength(
        r'diffSignal',
        0,
        true,
        0,
        true,
      );
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      diffSignalIsNotEmpty() {
    return QueryBuilder.apply(this, (query) {
      return query.listLength(
        r'diffSignal',
        0,
        false,
        999999,
        true,
      );
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      diffSignalLengthLessThan(
    int length, {
    bool include = false,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.listLength(
        r'diffSignal',
        0,
        true,
        length,
        include,
      );
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      diffSignalLengthGreaterThan(
    int length, {
    bool include = false,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.listLength(
        r'diffSignal',
        length,
        include,
        999999,
        true,
      );
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      diffSignalLengthBetween(
    int lower,
    int upper, {
    bool includeLower = true,
    bool includeUpper = true,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.listLength(
        r'diffSignal',
        lower,
        includeLower,
        upper,
        includeUpper,
      );
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      fingerprintElementEqualTo(
    double value, {
    double epsilon = Query.epsilon,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.equalTo(
        property: r'fingerprint',
        value: value,
        epsilon: epsilon,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      fingerprintElementGreaterThan(
    double value, {
    bool include = false,
    double epsilon = Query.epsilon,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.greaterThan(
        include: include,
        property: r'fingerprint',
        value: value,
        epsilon: epsilon,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      fingerprintElementLessThan(
    double value, {
    bool include = false,
    double epsilon = Query.epsilon,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.lessThan(
        include: include,
        property: r'fingerprint',
        value: value,
        epsilon: epsilon,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      fingerprintElementBetween(
    double lower,
    double upper, {
    bool includeLower = true,
    bool includeUpper = true,
    double epsilon = Query.epsilon,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.between(
        property: r'fingerprint',
        lower: lower,
        includeLower: includeLower,
        upper: upper,
        includeUpper: includeUpper,
        epsilon: epsilon,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      fingerprintLengthEqualTo(int length) {
    return QueryBuilder.apply(this, (query) {
      return query.listLength(
        r'fingerprint',
        length,
        true,
        length,
        true,
      );
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      fingerprintIsEmpty() {
    return QueryBuilder.apply(this, (query) {
      return query.listLength(
        r'fingerprint',
        0,
        true,
        0,
        true,
      );
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      fingerprintIsNotEmpty() {
    return QueryBuilder.apply(this, (query) {
      return query.listLength(
        r'fingerprint',
        0,
        false,
        999999,
        true,
      );
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      fingerprintLengthLessThan(
    int length, {
    bool include = false,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.listLength(
        r'fingerprint',
        0,
        true,
        length,
        include,
      );
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      fingerprintLengthGreaterThan(
    int length, {
    bool include = false,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.listLength(
        r'fingerprint',
        length,
        include,
        999999,
        true,
      );
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      fingerprintLengthBetween(
    int lower,
    int upper, {
    bool includeLower = true,
    bool includeUpper = true,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.listLength(
        r'fingerprint',
        lower,
        includeLower,
        upper,
        includeUpper,
      );
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      healthScoreEqualTo(int value) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.equalTo(
        property: r'healthScore',
        value: value,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      healthScoreGreaterThan(
    int value, {
    bool include = false,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.greaterThan(
        include: include,
        property: r'healthScore',
        value: value,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      healthScoreLessThan(
    int value, {
    bool include = false,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.lessThan(
        include: include,
        property: r'healthScore',
        value: value,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      healthScoreBetween(
    int lower,
    int upper, {
    bool includeLower = true,
    bool includeUpper = true,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.between(
        property: r'healthScore',
        lower: lower,
        includeLower: includeLower,
        upper: upper,
        includeUpper: includeUpper,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      idEqualTo(Id value) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.equalTo(
        property: r'id',
        value: value,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      idGreaterThan(
    Id value, {
    bool include = false,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.greaterThan(
        include: include,
        property: r'id',
        value: value,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      idLessThan(
    Id value, {
    bool include = false,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.lessThan(
        include: include,
        property: r'id',
        value: value,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      idBetween(
    Id lower,
    Id upper, {
    bool includeLower = true,
    bool includeUpper = true,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.between(
        property: r'id',
        lower: lower,
        includeLower: includeLower,
        upper: upper,
        includeUpper: includeUpper,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      isDeletedEqualTo(bool value) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.equalTo(
        property: r'isDeleted',
        value: value,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      isSyncedEqualTo(bool value) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.equalTo(
        property: r'isSynced',
        value: value,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      measuredAtEqualTo(DateTime value) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.equalTo(
        property: r'measuredAt',
        value: value,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      measuredAtGreaterThan(
    DateTime value, {
    bool include = false,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.greaterThan(
        include: include,
        property: r'measuredAt',
        value: value,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      measuredAtLessThan(
    DateTime value, {
    bool include = false,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.lessThan(
        include: include,
        property: r'measuredAt',
        value: value,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      measuredAtBetween(
    DateTime lower,
    DateTime upper, {
    bool includeLower = true,
    bool includeUpper = true,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.between(
        property: r'measuredAt',
        lower: lower,
        includeLower: includeLower,
        upper: upper,
        includeUpper: includeUpper,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      riskLabelEqualTo(
    String value, {
    bool caseSensitive = true,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.equalTo(
        property: r'riskLabel',
        value: value,
        caseSensitive: caseSensitive,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      riskLabelGreaterThan(
    String value, {
    bool include = false,
    bool caseSensitive = true,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.greaterThan(
        include: include,
        property: r'riskLabel',
        value: value,
        caseSensitive: caseSensitive,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      riskLabelLessThan(
    String value, {
    bool include = false,
    bool caseSensitive = true,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.lessThan(
        include: include,
        property: r'riskLabel',
        value: value,
        caseSensitive: caseSensitive,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      riskLabelBetween(
    String lower,
    String upper, {
    bool includeLower = true,
    bool includeUpper = true,
    bool caseSensitive = true,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.between(
        property: r'riskLabel',
        lower: lower,
        includeLower: includeLower,
        upper: upper,
        includeUpper: includeUpper,
        caseSensitive: caseSensitive,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      riskLabelStartsWith(
    String value, {
    bool caseSensitive = true,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.startsWith(
        property: r'riskLabel',
        value: value,
        caseSensitive: caseSensitive,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      riskLabelEndsWith(
    String value, {
    bool caseSensitive = true,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.endsWith(
        property: r'riskLabel',
        value: value,
        caseSensitive: caseSensitive,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      riskLabelContains(String value, {bool caseSensitive = true}) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.contains(
        property: r'riskLabel',
        value: value,
        caseSensitive: caseSensitive,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      riskLabelMatches(String pattern, {bool caseSensitive = true}) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.matches(
        property: r'riskLabel',
        wildcard: pattern,
        caseSensitive: caseSensitive,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      riskLabelIsEmpty() {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.equalTo(
        property: r'riskLabel',
        value: '',
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      riskLabelIsNotEmpty() {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.greaterThan(
        property: r'riskLabel',
        value: '',
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      serverIdIsNull() {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(const FilterCondition.isNull(
        property: r'serverId',
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      serverIdIsNotNull() {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(const FilterCondition.isNotNull(
        property: r'serverId',
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      serverIdEqualTo(
    String? value, {
    bool caseSensitive = true,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.equalTo(
        property: r'serverId',
        value: value,
        caseSensitive: caseSensitive,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      serverIdGreaterThan(
    String? value, {
    bool include = false,
    bool caseSensitive = true,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.greaterThan(
        include: include,
        property: r'serverId',
        value: value,
        caseSensitive: caseSensitive,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      serverIdLessThan(
    String? value, {
    bool include = false,
    bool caseSensitive = true,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.lessThan(
        include: include,
        property: r'serverId',
        value: value,
        caseSensitive: caseSensitive,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      serverIdBetween(
    String? lower,
    String? upper, {
    bool includeLower = true,
    bool includeUpper = true,
    bool caseSensitive = true,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.between(
        property: r'serverId',
        lower: lower,
        includeLower: includeLower,
        upper: upper,
        includeUpper: includeUpper,
        caseSensitive: caseSensitive,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      serverIdStartsWith(
    String value, {
    bool caseSensitive = true,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.startsWith(
        property: r'serverId',
        value: value,
        caseSensitive: caseSensitive,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      serverIdEndsWith(
    String value, {
    bool caseSensitive = true,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.endsWith(
        property: r'serverId',
        value: value,
        caseSensitive: caseSensitive,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      serverIdContains(String value, {bool caseSensitive = true}) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.contains(
        property: r'serverId',
        value: value,
        caseSensitive: caseSensitive,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      serverIdMatches(String pattern, {bool caseSensitive = true}) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.matches(
        property: r'serverId',
        wildcard: pattern,
        caseSensitive: caseSensitive,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      serverIdIsEmpty() {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.equalTo(
        property: r'serverId',
        value: '',
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      serverIdIsNotEmpty() {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.greaterThan(
        property: r'serverId',
        value: '',
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      syncedAtIsNull() {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(const FilterCondition.isNull(
        property: r'syncedAt',
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      syncedAtIsNotNull() {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(const FilterCondition.isNotNull(
        property: r'syncedAt',
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      syncedAtEqualTo(DateTime? value) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.equalTo(
        property: r'syncedAt',
        value: value,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      syncedAtGreaterThan(
    DateTime? value, {
    bool include = false,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.greaterThan(
        include: include,
        property: r'syncedAt',
        value: value,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      syncedAtLessThan(
    DateTime? value, {
    bool include = false,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.lessThan(
        include: include,
        property: r'syncedAt',
        value: value,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      syncedAtBetween(
    DateTime? lower,
    DateTime? upper, {
    bool includeLower = true,
    bool includeUpper = true,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.between(
        property: r'syncedAt',
        lower: lower,
        includeLower: includeLower,
        upper: upper,
        includeUpper: includeUpper,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      updatedAtEqualTo(DateTime value) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.equalTo(
        property: r'updatedAt',
        value: value,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      updatedAtGreaterThan(
    DateTime value, {
    bool include = false,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.greaterThan(
        include: include,
        property: r'updatedAt',
        value: value,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      updatedAtLessThan(
    DateTime value, {
    bool include = false,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.lessThan(
        include: include,
        property: r'updatedAt',
        value: value,
      ));
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterFilterCondition>
      updatedAtBetween(
    DateTime lower,
    DateTime upper, {
    bool includeLower = true,
    bool includeUpper = true,
  }) {
    return QueryBuilder.apply(this, (query) {
      return query.addFilterCondition(FilterCondition.between(
        property: r'updatedAt',
        lower: lower,
        includeLower: includeLower,
        upper: upper,
        includeUpper: includeUpper,
      ));
    });
  }
}

extension LocalMeasurementQueryObject
    on QueryBuilder<LocalMeasurement, LocalMeasurement, QFilterCondition> {}

extension LocalMeasurementQueryLinks
    on QueryBuilder<LocalMeasurement, LocalMeasurement, QFilterCondition> {}

extension LocalMeasurementQuerySortBy
    on QueryBuilder<LocalMeasurement, LocalMeasurement, QSortBy> {
  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      sortByClientLocalId() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'clientLocalId', Sort.asc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      sortByClientLocalIdDesc() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'clientLocalId', Sort.desc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      sortByDeviceMac() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'deviceMac', Sort.asc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      sortByDeviceMacDesc() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'deviceMac', Sort.desc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      sortByHealthScore() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'healthScore', Sort.asc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      sortByHealthScoreDesc() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'healthScore', Sort.desc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      sortByIsDeleted() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'isDeleted', Sort.asc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      sortByIsDeletedDesc() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'isDeleted', Sort.desc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      sortByIsSynced() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'isSynced', Sort.asc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      sortByIsSyncedDesc() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'isSynced', Sort.desc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      sortByMeasuredAt() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'measuredAt', Sort.asc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      sortByMeasuredAtDesc() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'measuredAt', Sort.desc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      sortByRiskLabel() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'riskLabel', Sort.asc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      sortByRiskLabelDesc() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'riskLabel', Sort.desc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      sortByServerId() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'serverId', Sort.asc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      sortByServerIdDesc() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'serverId', Sort.desc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      sortBySyncedAt() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'syncedAt', Sort.asc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      sortBySyncedAtDesc() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'syncedAt', Sort.desc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      sortByUpdatedAt() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'updatedAt', Sort.asc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      sortByUpdatedAtDesc() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'updatedAt', Sort.desc);
    });
  }
}

extension LocalMeasurementQuerySortThenBy
    on QueryBuilder<LocalMeasurement, LocalMeasurement, QSortThenBy> {
  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      thenByClientLocalId() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'clientLocalId', Sort.asc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      thenByClientLocalIdDesc() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'clientLocalId', Sort.desc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      thenByDeviceMac() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'deviceMac', Sort.asc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      thenByDeviceMacDesc() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'deviceMac', Sort.desc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      thenByHealthScore() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'healthScore', Sort.asc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      thenByHealthScoreDesc() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'healthScore', Sort.desc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy> thenById() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'id', Sort.asc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      thenByIdDesc() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'id', Sort.desc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      thenByIsDeleted() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'isDeleted', Sort.asc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      thenByIsDeletedDesc() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'isDeleted', Sort.desc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      thenByIsSynced() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'isSynced', Sort.asc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      thenByIsSyncedDesc() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'isSynced', Sort.desc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      thenByMeasuredAt() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'measuredAt', Sort.asc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      thenByMeasuredAtDesc() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'measuredAt', Sort.desc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      thenByRiskLabel() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'riskLabel', Sort.asc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      thenByRiskLabelDesc() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'riskLabel', Sort.desc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      thenByServerId() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'serverId', Sort.asc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      thenByServerIdDesc() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'serverId', Sort.desc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      thenBySyncedAt() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'syncedAt', Sort.asc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      thenBySyncedAtDesc() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'syncedAt', Sort.desc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      thenByUpdatedAt() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'updatedAt', Sort.asc);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QAfterSortBy>
      thenByUpdatedAtDesc() {
    return QueryBuilder.apply(this, (query) {
      return query.addSortBy(r'updatedAt', Sort.desc);
    });
  }
}

extension LocalMeasurementQueryWhereDistinct
    on QueryBuilder<LocalMeasurement, LocalMeasurement, QDistinct> {
  QueryBuilder<LocalMeasurement, LocalMeasurement, QDistinct>
      distinctByClientLocalId({bool caseSensitive = true}) {
    return QueryBuilder.apply(this, (query) {
      return query.addDistinctBy(r'clientLocalId',
          caseSensitive: caseSensitive);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QDistinct>
      distinctByDeviceMac({bool caseSensitive = true}) {
    return QueryBuilder.apply(this, (query) {
      return query.addDistinctBy(r'deviceMac', caseSensitive: caseSensitive);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QDistinct>
      distinctByDiffSignal() {
    return QueryBuilder.apply(this, (query) {
      return query.addDistinctBy(r'diffSignal');
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QDistinct>
      distinctByFingerprint() {
    return QueryBuilder.apply(this, (query) {
      return query.addDistinctBy(r'fingerprint');
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QDistinct>
      distinctByHealthScore() {
    return QueryBuilder.apply(this, (query) {
      return query.addDistinctBy(r'healthScore');
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QDistinct>
      distinctByIsDeleted() {
    return QueryBuilder.apply(this, (query) {
      return query.addDistinctBy(r'isDeleted');
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QDistinct>
      distinctByIsSynced() {
    return QueryBuilder.apply(this, (query) {
      return query.addDistinctBy(r'isSynced');
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QDistinct>
      distinctByMeasuredAt() {
    return QueryBuilder.apply(this, (query) {
      return query.addDistinctBy(r'measuredAt');
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QDistinct>
      distinctByRiskLabel({bool caseSensitive = true}) {
    return QueryBuilder.apply(this, (query) {
      return query.addDistinctBy(r'riskLabel', caseSensitive: caseSensitive);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QDistinct>
      distinctByServerId({bool caseSensitive = true}) {
    return QueryBuilder.apply(this, (query) {
      return query.addDistinctBy(r'serverId', caseSensitive: caseSensitive);
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QDistinct>
      distinctBySyncedAt() {
    return QueryBuilder.apply(this, (query) {
      return query.addDistinctBy(r'syncedAt');
    });
  }

  QueryBuilder<LocalMeasurement, LocalMeasurement, QDistinct>
      distinctByUpdatedAt() {
    return QueryBuilder.apply(this, (query) {
      return query.addDistinctBy(r'updatedAt');
    });
  }
}

extension LocalMeasurementQueryProperty
    on QueryBuilder<LocalMeasurement, LocalMeasurement, QQueryProperty> {
  QueryBuilder<LocalMeasurement, int, QQueryOperations> idProperty() {
    return QueryBuilder.apply(this, (query) {
      return query.addPropertyName(r'id');
    });
  }

  QueryBuilder<LocalMeasurement, String, QQueryOperations>
      clientLocalIdProperty() {
    return QueryBuilder.apply(this, (query) {
      return query.addPropertyName(r'clientLocalId');
    });
  }

  QueryBuilder<LocalMeasurement, String, QQueryOperations> deviceMacProperty() {
    return QueryBuilder.apply(this, (query) {
      return query.addPropertyName(r'deviceMac');
    });
  }

  QueryBuilder<LocalMeasurement, List<double>, QQueryOperations>
      diffSignalProperty() {
    return QueryBuilder.apply(this, (query) {
      return query.addPropertyName(r'diffSignal');
    });
  }

  QueryBuilder<LocalMeasurement, List<double>, QQueryOperations>
      fingerprintProperty() {
    return QueryBuilder.apply(this, (query) {
      return query.addPropertyName(r'fingerprint');
    });
  }

  QueryBuilder<LocalMeasurement, int, QQueryOperations> healthScoreProperty() {
    return QueryBuilder.apply(this, (query) {
      return query.addPropertyName(r'healthScore');
    });
  }

  QueryBuilder<LocalMeasurement, bool, QQueryOperations> isDeletedProperty() {
    return QueryBuilder.apply(this, (query) {
      return query.addPropertyName(r'isDeleted');
    });
  }

  QueryBuilder<LocalMeasurement, bool, QQueryOperations> isSyncedProperty() {
    return QueryBuilder.apply(this, (query) {
      return query.addPropertyName(r'isSynced');
    });
  }

  QueryBuilder<LocalMeasurement, DateTime, QQueryOperations>
      measuredAtProperty() {
    return QueryBuilder.apply(this, (query) {
      return query.addPropertyName(r'measuredAt');
    });
  }

  QueryBuilder<LocalMeasurement, String, QQueryOperations> riskLabelProperty() {
    return QueryBuilder.apply(this, (query) {
      return query.addPropertyName(r'riskLabel');
    });
  }

  QueryBuilder<LocalMeasurement, String?, QQueryOperations> serverIdProperty() {
    return QueryBuilder.apply(this, (query) {
      return query.addPropertyName(r'serverId');
    });
  }

  QueryBuilder<LocalMeasurement, DateTime?, QQueryOperations>
      syncedAtProperty() {
    return QueryBuilder.apply(this, (query) {
      return query.addPropertyName(r'syncedAt');
    });
  }

  QueryBuilder<LocalMeasurement, DateTime, QQueryOperations>
      updatedAtProperty() {
    return QueryBuilder.apply(this, (query) {
      return query.addPropertyName(r'updatedAt');
    });
  }
}
