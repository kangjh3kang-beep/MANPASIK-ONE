export interface DataPacket {
  id: number;
  timestamp: string;
  header: {
    version: string;
    type: string;
    sequence: number;
  };
  payload: {
    sensorId: string;
    value: number;
    unit: string;
  };
  checksum: string;
  prevHash: string;
  currentHash: string;
  isValid: boolean;
  validationError?: string;
}

export class PacketGenerator {
  private sequence: number = 0;
  private previousHash: string = '0000000000000000';

  private calculateHash(data: string): string {
    let hash = 0;
    for (let i = 0; i < data.length; i++) {
      const char = data.charCodeAt(i);
      hash = ((hash << 5) - hash) + char;
      hash = hash & hash;
    }
    return Math.abs(hash).toString(16).padStart(16, '0');
  }

  private calculateChecksum(header: any, payload: any): string {
    const data = JSON.stringify({ header, payload });
    return this.calculateHash(data).slice(0, 8);
  }

  generatePacket(injectError: boolean = false): DataPacket {
    this.sequence++;

    const header = {
      version: '1.0',
      type: 'SENSOR_DATA',
      sequence: this.sequence,
    };

    const payload = {
      sensorId: `SENSOR_${Math.floor(Math.random() * 5) + 1}`,
      value: parseFloat((20 + Math.random() * 30).toFixed(2)),
      unit: 'celsius',
    };

    const correctChecksum = this.calculateChecksum(header, payload);
    const checksum = injectError
      ? correctChecksum.slice(0, -2) + 'FF'
      : correctChecksum;

    const hashInput = JSON.stringify({ header, payload, checksum, prevHash: this.previousHash });
    const currentHash = this.calculateHash(hashInput);

    const packet: DataPacket = {
      id: this.sequence,
      timestamp: new Date().toISOString(),
      header,
      payload,
      checksum,
      prevHash: this.previousHash,
      currentHash,
      isValid: true,
      validationError: undefined,
    };

    const recalculatedChecksum = this.calculateChecksum(header, payload);
    if (checksum !== recalculatedChecksum) {
      packet.isValid = false;
      packet.validationError = 'Checksum mismatch';
    }

    this.previousHash = currentHash;

    return packet;
  }

  validatePacket(packet: DataPacket, previousPacket: DataPacket | null): {
    isValid: boolean;
    errors: string[];
  } {
    const errors: string[] = [];

    const recalculatedChecksum = this.calculateChecksum(packet.header, packet.payload);
    if (packet.checksum !== recalculatedChecksum) {
      errors.push(`Checksum verification failed: expected ${recalculatedChecksum}, got ${packet.checksum}`);
    }

    if (previousPacket && packet.prevHash !== previousPacket.currentHash) {
      errors.push(`Hash chain broken: expected ${previousPacket.currentHash}, got ${packet.prevHash}`);
    }

    if (packet.header.sequence !== this.sequence) {
      errors.push(`Sequence number mismatch: expected ${this.sequence}, got ${packet.header.sequence}`);
    }

    return {
      isValid: errors.length === 0,
      errors,
    };
  }

  reset(): void {
    this.sequence = 0;
    this.previousHash = '0000000000000000';
  }
}
