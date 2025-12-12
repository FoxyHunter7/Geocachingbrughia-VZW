<?php

namespace Database\Seeders;

use App\Models\Social;
use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\DB;

class SocialsSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        DB::table('socials')->delete();

        $data = ReadFromCsv::getDataFromCsv(storage_path('app/data/csv/socials.csv'), ';');

        $model = new Social();
        foreach ($data as $row) {
            $model->create($row);
        }
    }
}
