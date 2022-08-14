import {inject, TestBed} from '@angular/core/testing';

import {ResortService} from './resort.service';
import {count, map, reduce, tap, toArray} from "rxjs";
import {HttpClientTestingModule, HttpTestingController} from "@angular/common/http/testing";
import {ResortsEntity} from "@ui/data";

describe('ResortService', () => {

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [ResortService]
    });
  });

  it('should be created', inject([ResortService], (service: ResortService) => {
    expect(service).toBeTruthy();
  }));

  it('should retrieve resorts', inject(
    [HttpTestingController, ResortService],
    (httpMock: HttpTestingController, service: ResortService) => {
      const mockResorts = [
        {code: "akl", name: "Animal Kingdom Lodge"} as ResortsEntity,
        {code: "ssr", name: "Saratoga Springs Resport"} as ResortsEntity,
        // {code: "blt", name: "Bay Lake Tower"} as ResortsEntity,
        {code: "bcv", name: "Beach Club Villas"} as ResortsEntity,
        {code: "bwv", name: "Boardwalk Villas"} as ResortsEntity,
      ]
      service.$loadResorts()
        .pipe(toArray())
        .subscribe((result) => {
          expect(result.length).toBe(5)
          expect(result).toBe(mockResorts)
        });
      const mockReq = httpMock.expectOne('http://localhost:8080/resort');
      expect(mockReq.cancelled).toBeFalsy();
      expect(mockReq.request.method).toBe('GET')
      expect(mockReq.request.responseType).toEqual('json');
      mockReq.flush(mockResorts);
      httpMock.verify();
    }))
});
